package server

import (
	"context"
	"errors"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/server/object"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/maps"
	"sync"
)

var (
	inst = &service{
		servers: make(map[string]*object.ServerInfo),
		groups:  make(map[string]*object.ServerGroup),
	}

	collectionServers *mongo.Collection
	collectionGroups  *mongo.Collection
)

// Service is a service that provides server information.
// It is used to cache server information and look up servers by their ID or port.
// It is also used to cache server groups and look up server groups by their ID.
// The service is thread-safe.
type service struct {
	servers   map[string]*object.ServerInfo
	serversMu sync.Mutex

	groups   map[string]*object.ServerGroup
	groupsMu sync.Mutex
}

// LookupById looks up a server by its ID.
func (s *service) LookupById(id string) *object.ServerInfo {
	s.serversMu.Lock()
	defer s.serversMu.Unlock()

	server, ok := s.servers[id]
	if !ok {
		return nil
	}

	return server
}

// LookupByPort looks up a server by its port.
func (s *service) LookupByPort(port int64) *object.ServerInfo {
	s.serversMu.Lock()
	defer s.serversMu.Unlock()

	for _, server := range s.servers {
		if server.Port() == port {
			return server
		}
	}

	return nil
}

// CacheServer caches a server in the service.
func (s *service) CacheServer(server *object.ServerInfo) {
	s.serversMu.Lock()
	s.servers[server.Id()] = server
	s.serversMu.Unlock()
}

// DestroyServer removes a server from the cache.
func (s *service) DestroyServer(id string) {
	s.serversMu.Lock()
	delete(s.servers, id)
	s.serversMu.Unlock()
}

func (s *service) Servers() []*object.ServerInfo {
	s.serversMu.Lock()
	defer s.serversMu.Unlock()

	return maps.Values(s.servers)
}

// LookupGroupById looks up a server group by its ID.
func (s *service) LookupGroupById(id string) *object.ServerGroup {
	s.groupsMu.Lock()
	defer s.groupsMu.Unlock()

	group, ok := s.groups[id]
	if !ok {
		return nil
	}

	return group
}

// CacheGroup caches a server group in the service.
func (s *service) CacheGroup(group *object.ServerGroup) {
	s.groupsMu.Lock()
	s.groups[group.Id()] = group
	s.groupsMu.Unlock()
}

// DestroyGroup removes a server group from the cache.
func (s *service) DestroyGroup(id string) {
	s.groupsMu.Lock()
	delete(s.groups, id)
	s.groupsMu.Unlock()
}

func (s *service) Groups() []*object.ServerGroup {
	s.groupsMu.Lock()
	defer s.groupsMu.Unlock()

	return maps.Values(s.groups)
}

// Service returns the server service instance.
func Service() *service {
	return inst
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

			i, err := object.Unmarshal(result)
			if err != nil {
				return errors.Join(errors.New("failed to unmarshal server"), err)
			}

			inst.CacheServer(i)
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

			g := &object.ServerGroup{}
			if err := g.Unmarshal(result); err != nil {
				return errors.Join(errors.New("failed to unmarshal server group"), err)
			}

			inst.CacheGroup(g)
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
