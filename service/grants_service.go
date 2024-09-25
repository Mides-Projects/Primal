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
	accountsMu sync.RWMutex
	accounts   map[string]*grantsx.Tracker

	col *mongo.Collection
}

// Lookup retrieves a Tracker from the cache by its ID.
func (s *GrantsService) Lookup(id string) *grantsx.Tracker {
	s.accountsMu.RLock()
	defer s.accountsMu.RUnlock()

	return s.accounts[id]
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

	acc, err := accountService.UnsafeLookupById(id)
	if err != nil {
		return nil, err
	}

	ga := grantsx.EmptyGrantsAccount(acc)
	for cur.Next(context.Background()) {
		var body map[string]interface{}
		if err := cur.Decode(&body); err != nil {
			return nil, err
		}

		g := &grantsx.Grant{}
		if err := g.Unmarshal(body); err != nil {
			return nil, err
		}

		if g.Expired() {
			ga.AddExpiredGrant(g)
		} else {
			ga.AddActiveGrant(g)
		}
	}

	s.Cache(ga)

	return ga, nil
}

// HighestGroupBy retrieves the highest group by its ID.
func (s *GrantsService) HighestGroupBy(ga *grantsx.Tracker) *grantsx.Group {
	var highest *grantsx.Group
	for _, gr := range ga.ActiveGrants() {
		if gr.Identifier().Key() != "group" {
			continue
		}

		g := groupsService.LookupById(gr.Identifier().Value())
		if g == nil {
			continue
		}

		if highest == nil || g.Weight() > highest.Weight() {
			highest = g
		}
	}

	return highest
}

// Cache caches a Tracker.
func (s *GrantsService) Cache(ga *grantsx.Tracker) {
	s.accountsMu.Lock()
	s.accounts[ga.Account().Id()] = ga
	s.accountsMu.Unlock()
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
	accounts: make(map[string]*grantsx.Tracker),
}
