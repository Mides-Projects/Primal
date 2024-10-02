package service

import (
	"context"
	"errors"
	"github.com/holypvp/primal/model/grantsx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

type GrantsService struct {
	trackersMu sync.RWMutex
	trackers   map[string]*grantsx.Tracker

	col *mongo.Collection
}

// Lookup retrieves a Tracker from the cache by its ID.
func (s *GrantsService) Lookup(id string) *grantsx.Tracker {
	s.trackersMu.RLock()
	defer s.trackersMu.RUnlock()

	return s.trackers[id]
}

// UnsafeLookup retrieves a Tracker from the cache by its ID.
// If the account is not found in the cache, it will be fetched from the database.
// This method is not thread-safe and should be used with caution in goroutines.
func (s *GrantsService) UnsafeLookup(id string) (*grantsx.Tracker, error) {
	if ga := s.Lookup(id); ga != nil {
		return ga, nil
	}

	if s.col == nil {
		return nil, errors.New("service not hooked to the database")
	}

	cur, err := s.col.Find(context.Background(), bson.M{"source_id": id})
	if err != nil {
		return nil, err
	}

	ga := grantsx.EmptyTracker()
	for cur.Next(context.Background()) {
		var body map[string]interface{}
		if err = cur.Decode(&body); err != nil {
			return nil, err
		}

		g := &grantsx.Grant{}
		if err = g.Unmarshal(body); err != nil {
			return nil, err
		}

		if g.Expired() {
			ga.AddExpiredGrant(g)
		} else {
			ga.AddActiveGrant(g)
		}
	}

	s.Cache(id, ga)

	return ga, nil
}

// Cache caches a Tracker.
func (s *GrantsService) Cache(id string, t *grantsx.Tracker) {
	s.trackersMu.Lock()
	s.trackers[id] = t
	s.trackersMu.Unlock()
}

// Save saves a grant.
func (s *GrantsService) Save(srcId string, g *grantsx.Grant) error {
	if s.col == nil {
		return errors.New("service not hooked to the database")
	}

	body := g.Marshal()
	body["source_id"] = srcId

	_, err := s.col.UpdateOne(
		context.TODO(),
		bson.M{"_id": g.Id()},
		bson.M{"_set": body},
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *GrantsService) Hook(db *mongo.Database) error {
	if s.col != nil {
		return errors.New("an instance of GrantsService already exists")
	}

	s.col = db.Collection("grantsx")

	return nil
}

func Grants() *GrantsService {
	return grantsService
}

var grantsService = &GrantsService{
	trackers: make(map[string]*grantsx.Tracker),
}
