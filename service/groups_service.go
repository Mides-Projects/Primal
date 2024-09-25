package service

import (
	"context"
	"errors"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/grantsx/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/maps"
	"strings"
	"sync"
)

type GroupsService struct {
	groupsMu sync.RWMutex
	groups   map[string]*model.Group

	col *mongo.Collection
}

// All returns all bgroups.
func (s *GroupsService) All() []*model.Group {
	s.groupsMu.RLock()
	defer s.groupsMu.RUnlock()

	return maps.Values(s.groups)
}

// LookupById looks up a group by its ID.
func (s *GroupsService) LookupById(id string) *model.Group {
	s.groupsMu.RLock()
	defer s.groupsMu.RUnlock()

	g, ok := s.groups[id]
	if !ok {
		return nil
	}

	return g
}

// LookupByName looks up a group by its name.
func (s *GroupsService) LookupByName(name string) *model.Group {
	s.groupsMu.RLock()
	defer s.groupsMu.RUnlock()

	name = strings.ToLower(name)
	for _, g := range s.groups {
		if strings.ToLower(g.Name()) != name {
			continue
		}

		return g
	}

	return nil
}

// Cache caches a group.
func (s *GroupsService) Cache(g *model.Group) {
	s.groupsMu.Lock()
	s.groups[g.Id()] = g
	s.groupsMu.Unlock()
}

func (s *GroupsService) Save(g *model.Group) error {
	if s.col == nil {
		return errors.New("service not hooked to the database")
	}

	res, err := s.col.UpdateOne(
		context.TODO(),
		bson.D{{"_id", g.Id()}},
		bson.D{{"$set", g.Marshal()}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return err
	}

	if res.UpsertedCount == 0 {
		common.Log.Printf("Group %s was updated", g.Id())
	} else {
		common.Log.Printf("Group %s was inserted", g.Id())
	}

	return nil
}

// Hook initializes the service with the database.
func (s *GroupsService) Hook(db *mongo.Database) error {
	if s.col != nil {
		return errors.New("an instance of Service for 'Groups' already exists")
	}

	s.col = db.Collection("bgroups")

	// Load bgroups from the database
	cur, err := s.col.Find(context.Background(), bson.D{})
	if err != nil {
		return err
	}

	for cur.Next(context.Background()) {
		var body map[string]interface{}
		if err = cur.Decode(&body); err != nil {
			return err
		}

		g := &model.Group{}
		if err := g.Unmarshal(body); err != nil {
			return err
		}

		s.Cache(g)
	}

	common.Log.Printf("Successfully loaded %d groups", len(s.groups))

	return nil
}

func Groups() *GroupsService {
	return groupsService
}

var groupsService = &GroupsService{
	groups: make(map[string]*model.Group),
}
