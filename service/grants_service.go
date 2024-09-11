package service

import (
	"context"
	"errors"
	"github.com/holypvp/primal/grantsx/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

type GrantsService struct {
	accountsMu sync.RWMutex
	accounts   map[string]*model.GrantsAccount

	col *mongo.Collection
}

// LookupAtCache retrieves a GrantsAccount from the cache by its ID.
func (s *GrantsService) LookupAtCache(id string) *model.GrantsAccount {
	s.accountsMu.RLock()
	defer s.accountsMu.RUnlock()

	return s.accounts[id]
}

func (s *GrantsService) Lookup(id string) (*model.GrantsAccount, error) {
	if ga := s.LookupAtCache(id); ga != nil {
		return ga, nil
	}

	if s.col == nil {
		return nil, errors.New("service not hooked to the database")
	}

	cur, err := s.col.Find(context.Background(), bson.M{"source_id": id})
	if err != nil {
		return nil, err
	}

	acc, err := accountService.Fetch(id)
	if err != nil {
		return nil, err
	}

	ga := model.EmptyGrantsAccount(acc)
	for cur.Next(context.Background()) {
		var body map[string]interface{}
		if err := cur.Decode(&body); err != nil {
			return nil, err
		}

		g := &model.Grant{}
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
func (s *GrantsService) HighestGroupBy(ga *model.GrantsAccount) *model.Group {
	var highest *model.Group
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

// Cache caches a GrantsAccount.
func (s *GrantsService) Cache(ga *model.GrantsAccount) {
	s.accountsMu.Lock()
	s.accounts[ga.Account().Id()] = ga
	s.accountsMu.Unlock()
}

// Save saves a grant.
func (s *GrantsService) Save(srcId string, g *model.Grant) error {
	if s.col == nil {
		return errors.New("service not hooked to the database")
	}

	body := g.Marshal()
	body["source_xuid"] = srcId

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

	s.col = db.Collection("grants")

	return nil
}

func Grants() *GrantsService {
	return grantsService
}

var grantsService = &GrantsService{
	accounts: make(map[string]*model.GrantsAccount),
}
