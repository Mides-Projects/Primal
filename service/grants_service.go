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

	var body map[string]interface{}
	if err := s.col.FindOne(context.Background(), bson.M{"_id": id}).Decode(&body); err != nil {
		return nil, err
	}

	ga := &model.GrantsAccount{}
	if err := ga.Unmarshal(body); err != nil {
		return nil, err
	}

	return ga, nil
}

func (s *GrantsService) Cache(ga *model.GrantsAccount) {
	s.accountsMu.Lock()
	s.accounts[ga.Account().Id()] = ga
	s.accountsMu.Unlock()
}

func (s *GrantsService) Hook(db *mongo.Database) error {
	if s.col != nil {
		return errors.New("an instance of GrantsService already exists")
	}

	s.col = db.Collection("grants")

	cur, err := s.col.Find(context.Background(), bson.D{})
	if err != nil {
		return err
	}

	for cur.Next(context.Background()) {
		var body map[string]interface{}
		if err := cur.Decode(&body); err != nil {
			return err
		}

		ga := &model.GrantsAccount{}
		if err := ga.Unmarshal(body); err != nil {
			return err
		}

		s.Cache(ga)
	}

	return nil
}

func Grants() *GrantsService {
	return grantsService
}

var grantsService = &GrantsService{
	accounts: make(map[string]*model.GrantsAccount),
}
