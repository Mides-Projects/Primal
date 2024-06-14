package server

import (
	"context"
	"github.com/holypvp/primal/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
	"sync"
)

var (
	serversMu sync.Mutex
	groupsMu  sync.Mutex
	portMu    sync.Mutex

	serversMap    = map[string]*ServerInfo{}
	serversByPort = map[int64]string{}

	groupsMap = map[string]*ServerGroup{}

	collectionServers *mongo.Collection
	collectionGroups  *mongo.Collection

	instance *ServerService
)

type ServerService struct{}

func (service *ServerService) LookupById(id string) *ServerInfo {
	serversMu.Lock()
	defer serversMu.Unlock()

	serverInfo, ok := serversMap[strings.ToLower(id)]
	if !ok {
		return nil
	}

	return serverInfo
}

func (service *ServerService) LookupByPort(port int64) *ServerInfo {
	portMu.Lock()
	defer portMu.Unlock()

	id, ok := serversByPort[port]
	if !ok {
		return nil
	}

	return service.LookupById(id)
}

func (service *ServerService) AppendServer(serverInfo *ServerInfo) {
	serversMu.Lock()
	serversMap[strings.ToLower(serverInfo.Id())] = serverInfo
	serversMu.Unlock()

	portMu.Lock()
	serversByPort[serverInfo.Port()] = serverInfo.Id()
	portMu.Unlock()
}

func (service *ServerService) DestroyServer(serverId string) {
	serversMu.Lock()
	defer serversMu.Unlock()

	serverInfo, ok := serversMap[strings.ToLower(serverId)]
	if !ok {
		return
	}

	portMu.Lock()
	defer portMu.Unlock()

	delete(serversMap, strings.ToLower(serverId))
	delete(serversByPort, serverInfo.Port())
}

func (service *ServerService) Servers() []*ServerInfo {
	serversMu.Lock()
	defer serversMu.Unlock()

	values := make([]*ServerInfo, 0, len(serversMap))
	for _, serverInfo := range serversMap {
		values = append(values, serverInfo)
	}

	return values
}

func (service *ServerService) LookupGroup(name string) *ServerGroup {
	groupsMu.Lock()
	defer groupsMu.Unlock()

	group, ok := groupsMap[strings.ToLower(name)]
	if !ok {
		return nil
	}

	return group
}

func (service *ServerService) AppendGroup(group *ServerGroup) {
	groupsMu.Lock()
	groupsMap[strings.ToLower(group.Id())] = group
	groupsMu.Unlock()
}

func (service *ServerService) DestroyGroup(name string) {
	groupsMu.Lock()
	defer groupsMu.Unlock()

	delete(groupsMap, strings.ToLower(name))
}

func (service *ServerService) Groups() []*ServerGroup {
	groupsMu.Lock()
	defer groupsMu.Unlock()

	values := make([]*ServerGroup, 0, len(groupsMap))
	for _, group := range groupsMap {
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
	if collectionServers == nil {
		log.Fatal("Servers collection is nil")
	}

	result, err := collectionServers.UpdateOne(
		context.TODO(),
		bson.D{{"_id", infoModel.Id}},
		bson.D{{"$set", infoModel}},
		options.Update().SetUpsert(true),
	)
	// _, err := collectionServers.InsertOne(context.TODO(), infoModel)
	if err != nil {
		log.Fatal(err)

		return
	}

	if result.UpsertedCount > 0 {
		log.Printf("Server %s was inserted", infoModel.Id)
	} else {
		log.Printf("Server %s was updated", infoModel.Id)
	}
}

func (service *ServerService) LoadGroups(database *mongo.Database) {
	if collectionGroups != nil {
		log.Fatal("Groups collection is already set")
	}

	collectionGroups = database.Collection("serverGroups")
	cursor, err := collectionGroups.Find(context.TODO(), bson.D{{}})
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

	log.Printf("Successfully loaded %d server groups", len(groupsMap))
}

func (service *ServerService) LoadServers(database *mongo.Database) {
	collectionServers = database.Collection("servers")
	if collectionServers == nil {
		log.Fatal("Servers collection is nil")
	}

	cursor, err := collectionServers.Find(context.TODO(), bson.D{{}})
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

	log.Printf("Successfully loaded %d servers", len(serversMap))
}
