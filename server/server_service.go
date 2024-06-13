package server

import (
	"context"
	"github.com/holypvp/primal/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strings"
	"sync"
)

var (
	serversMutex sync.Mutex
	groupsMutex  sync.Mutex
	portMutex    sync.Mutex

	servers       = map[string]*ServerInfo{}
	serversByPort = map[int64]string{}

	groups = map[string]*ServerGroup{}

	serversCollection *mongo.Collection
	groupsCollection  *mongo.Collection

	instance *ServerService
)

type ServerService struct{}

func (service *ServerService) LookupById(id string) *ServerInfo {
	serversMutex.Lock()
	defer serversMutex.Unlock()

	serverInfo, ok := servers[strings.ToLower(id)]
	if !ok {
		return nil
	}

	return serverInfo
}

func (service *ServerService) LookupByPort(port int64) *ServerInfo {
	portMutex.Lock()
	defer portMutex.Unlock()

	id, ok := serversByPort[port]
	if !ok {
		return nil
	}

	return service.LookupById(id)
}

func (service *ServerService) AppendServer(serverInfo *ServerInfo) {
	serversMutex.Lock()
	servers[strings.ToLower(serverInfo.Id())] = serverInfo
	serversMutex.Unlock()

	portMutex.Lock()
	serversByPort[serverInfo.Port()] = serverInfo.Id()
	portMutex.Unlock()
}

func (service *ServerService) DestroyServer(serverId string) {
	serversMutex.Lock()
	defer serversMutex.Unlock()

	serverInfo, ok := servers[strings.ToLower(serverId)]
	if !ok {
		return
	}

	portMutex.Lock()
	defer portMutex.Unlock()

	delete(servers, strings.ToLower(serverId))
	delete(serversByPort, serverInfo.Port())
}

func (service *ServerService) Servers() []*ServerInfo {
	serversMutex.Lock()
	defer serversMutex.Unlock()

	values := make([]*ServerInfo, 0, len(servers))
	for _, serverInfo := range servers {
		values = append(values, serverInfo)
	}

	return values
}

func (service *ServerService) LookupGroup(name string) *ServerGroup {
	groupsMutex.Lock()
	defer groupsMutex.Unlock()

	group, ok := groups[strings.ToLower(name)]
	if !ok {
		return nil
	}

	return group
}

func (service *ServerService) AppendGroup(group *ServerGroup) {
	groupsMutex.Lock()
	groups[strings.ToLower(group.Id())] = group
	groupsMutex.Unlock()
}

func (service *ServerService) DestroyGroup(name string) {
	groupsMutex.Lock()
	defer groupsMutex.Unlock()

	delete(groups, strings.ToLower(name))
}

func (service *ServerService) Groups() []*ServerGroup {
	groupsMutex.Lock()
	defer groupsMutex.Unlock()

	values := make([]*ServerGroup, 0, len(groups))
	for _, group := range groups {
		values = append(values, group)
	}

	return values
}

func Service() *ServerService {
	if instance == nil {
		instance = &ServerService{}
	}

	return instance
}

func SaveModel(infoModel model.ServerInfoModel) {
	if serversCollection == nil {
		log.Fatal("Servers collection is nil")
	}

	_, err := serversCollection.InsertOne(context.TODO(), infoModel)
	if err == nil {
		return
	}

	log.Fatal(err)
}

func (service *ServerService) LoadGroups(database *mongo.Database) {
	if groupsCollection != nil {
		log.Fatal("Groups collection is already set")
	}

	groupsCollection := database.Collection("serverGroups")

	cursor, err := groupsCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		panic(err)
	}

	for cursor.Next(context.Background()) {
		var result = &model.ServerGroupModel{}
		err := cursor.Decode(result)
		if err != nil {
			panic(err)

			return
		}

		service.AppendGroup(&ServerGroup{
			id:                    result.Id,
			metadata:              result.Metadata,
			announcements:         result.Announcements,
			announcementsInterval: result.AnnouncementsInterval,
			fallbackServerId:      result.FallbackServerId,
		})
	}

	if err := cursor.Err(); err != nil {
		panic(err)

		return
	}

	err = cursor.Close(context.TODO())
	if err != nil {
		panic(err)
	}

	log.Printf("Successfully loaded %d server groups\n\n", len(groups))
}

func (service *ServerService) LoadServers(database *mongo.Database) {
	serversCollection := database.Collection("servers")
	cursor, err := serversCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		panic(err)
	}

	for cursor.Next(context.Background()) {
		var result = &model.ServerInfoModel{}
		err := cursor.Decode(result)
		if err != nil {
			log.Println("Failed to decode server info model: ", err)

			continue
		}

		service.AppendServer(&ServerInfo{
			id:     result.Id,
			port:   result.Port,
			groups: result.Groups,

			maxSlots:    result.MaxSlots,
			heartbeat:   result.Heartbeat,
			bungeeCord:  result.BungeeCord,
			onlineMode:  result.OnlineMode,
			initialTime: result.InitialTime,
		})
	}

	if err := cursor.Err(); err != nil {
		panic(err)

		return
	}

	err = cursor.Close(context.TODO())
	if err != nil {
		panic(err)
	}

	log.Printf("Successfully loaded %d servers\n", len(servers))
}
