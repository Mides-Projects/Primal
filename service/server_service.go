package service

import (
	"context"
	"errors"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/maps"
	"sync"
)

// ServerService is a service for managing servers and server groups.
// It is used to cache server information and look up servers by their ID or port.
// It is also used to cache server groups and look up server groups by their ID.
// The ServerService is thread-safe.
type ServerService struct {
	servers   map[string]*model.ServerInfo
	serversMu sync.Mutex

	groups   map[string]*model.ServerGroup
	groupsMu sync.Mutex
}

// LookupById looks up a server by its ID.
func (s *ServerService) LookupById(id string) *model.ServerInfo {
	s.serversMu.Lock()
	defer s.serversMu.Unlock()

	server, ok := s.servers[id]
	if !ok {
		return nil
	}

	return server
}

// LookupByPort looks up a server by its port.
func (s *ServerService) LookupByPort(port int64) *model.ServerInfo {
	s.serversMu.Lock()
	defer s.serversMu.Unlock()

	for _, server := range s.servers {
		if server.Port() == port {
			return server
		}
	}

	return nil
}

// CacheServer caches a server in the ServerService.
func (s *ServerService) CacheServer(server *model.ServerInfo) {
	s.serversMu.Lock()
	s.servers[server.Id()] = server
	s.serversMu.Unlock()
}

// DestroyServer removes a server from the cache.
func (s *ServerService) DestroyServer(id string) {
	s.serversMu.Lock()
	delete(s.servers, id)
	s.serversMu.Unlock()
}

func (s *ServerService) Servers() []*model.ServerInfo {
	s.serversMu.Lock()
	defer s.serversMu.Unlock()

	return maps.Values(s.servers)
}

// LookupGroupById looks up a server group by its ID.
func (s *ServerService) LookupGroupById(id string) *model.ServerGroup {
	s.groupsMu.Lock()
	defer s.groupsMu.Unlock()

	group, ok := s.groups[id]
	if !ok {
		return nil
	}

	return group
}

// CacheGroup caches a server group in the ServerService.
func (s *ServerService) CacheGroup(group *model.ServerGroup) {
	s.groupsMu.Lock()
	s.groups[group.Id()] = group
	s.groupsMu.Unlock()
}

// DestroyGroup removes a server group from the cache.
func (s *ServerService) DestroyGroup(id string) {
	s.groupsMu.Lock()
	delete(s.groups, id)
	s.groupsMu.Unlock()
}

func (s *ServerService) Groups() []*model.ServerGroup {
	s.groupsMu.Lock()
	defer s.groupsMu.Unlock()

	return maps.Values(s.groups)
}

func SaveModel(id string, m map[string]interface{}) error {
	if collectionServers == nil {
		return errors.New("servers collection is not set")
	}

	result, err := collectionServers.UpdateOne(
		context.TODO(),
		bson.D{{"_id", id}},
		bson.D{{"$set", m}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return errors.Join(errors.New("failed to update server"), err)
	}

	if result.UpsertedCount > 0 {
		common.Log.Printf("Server %s was inserted", id)
	} else {
		common.Log.Printf("Server %s was updated", id)
	}

	return nil
}

func LoadServers(db *mongo.Database) error {
	if collectionServers != nil {
		return errors.New("servers collection is already set")
	}

	if c := db.Collection("servers"); c != nil {
		cursor, err := c.Find(context.TODO(), bson.D{{}})
		if err != nil {
			return errors.Join(errors.New("failed to load servers"), err)
		}

		for cursor.Next(context.Background()) {
			var result map[string]interface{}
			if err = cursor.Decode(&result); err != nil {
				return errors.Join(errors.New("failed to decode server"), err)
			}

			i, err := model.Unmarshal(result)
			if err != nil {
				return errors.Join(errors.New("failed to unmarshal server"), err)
			}

			serverService.CacheServer(i)
		}

		if err := cursor.Err(); err != nil {
			return errors.Join(errors.New("cursor error"), err)
		}

		if err = cursor.Close(context.TODO()); err != nil {
			return errors.Join(errors.New("failed to close cursor"), err)
		}

		collectionServers = c

		return nil
	}

	return errors.New("servers collection is nil")
}

func LoadGroups(database *mongo.Database) error {
	if collectionGroups != nil {
		return errors.New("server groups collection is already set")
	}

	if c := database.Collection("serverGroups"); c != nil {
		cursor, err := c.Find(context.TODO(), bson.D{{}})
		if err != nil {
			return errors.Join(errors.New("failed to load server groups"), err)
		}

		for cursor.Next(context.Background()) {
			var result map[string]interface{}
			if err = cursor.Decode(&result); err != nil {
				return errors.Join(errors.New("failed to decode server group"), err)
			}

			g := &model.ServerGroup{}
			if err := g.Unmarshal(result); err != nil {
				return errors.Join(errors.New("failed to unmarshal server group"), err)
			}

			serverService.CacheGroup(g)
		}

		if err := cursor.Err(); err != nil {
			return errors.Join(errors.New("cursor error"), err)
		}

		if err = cursor.Close(context.TODO()); err != nil {
			return errors.Join(errors.New("failed to close cursor"), err)
		}

		collectionGroups = c

		return nil
	}

	return errors.New("server groups collection is nil")
}

// Server returns the server ServerService instance.
func Server() *ServerService {
	return serverService
}

var (
	serverService = &ServerService{
		servers: make(map[string]*model.ServerInfo),
		groups:  make(map[string]*model.ServerGroup),
	}

	collectionServers *mongo.Collection
	collectionGroups  *mongo.Collection
)
