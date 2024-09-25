package service

import (
	"context"
	"errors"
	"github.com/holypvp/primal/account"
	"github.com/holypvp/primal/common"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"sync"
	"time"
)

type AccountService struct {
	col *mongo.Collection

	accountsMu sync.RWMutex
	accounts   map[string]*account.Account

	accountsIdMu sync.RWMutex
	accountsId   map[string]string
}

// LookupById retrieves an account by its ID. It's safe to use this method because
// it's only reading the map.
func (s *AccountService) LookupById(id string) *account.Account {
	s.accountsMu.RLock()
	defer s.accountsMu.RUnlock()

	return s.accounts[id]
}

// LookupByName retrieves an account by its name. It's safe to use this method because
// it's only reading the map.
func (s *AccountService) LookupByName(name string) *account.Account {
	s.accountsIdMu.RLock()
	defer s.accountsIdMu.RUnlock()

	if id, ok := s.accountsId[strings.ToLower(name)]; ok {
		return s.LookupById(id)
	}

	return nil
}

// UnsafeLookupById retrieves an account by its ID. It's unsafe to use this method because
// it's reading from the Redis database.
func (s *AccountService) UnsafeLookupById(id string) (*account.Account, error) {
	if acc := s.LookupById(id); acc != nil {
		return acc, nil
	}

	// These values are only cached for 72 hours, after that they are removed from redis
	// but still available in our mongo database.
	// If the account was fetch from database, it will be cached into redis to prevent further database calls in the next 72 hours.
	var (
		acc *account.Account
		err error
	)
	if acc, err = s.lookupAtRedis("primal%ids:", id); err != nil {
		return nil, err
	} else if acc == nil {
		acc, err = s.lookupAtMongo("_id", id)
	}

	if err != nil {
		return nil, err
	} else if acc == nil {
		return nil, nil
	}

	s.Cache(acc)

	return acc, nil

	// val, err := common.RedisClient.Get(context.Background(), "primal%ids:"+id).Result()
	// if errors.Is(err, redis.Nil) {
	// 	return nil, nil
	// } else if err != nil {
	// 	return nil, err
	// } else if val == "" {
	// 	return nil, errors.New("empty value")
	// } else {
	// 	acc := &account.Account{}
	// 	if err = acc.UnmarshalString(val); err != nil {
	// 		return nil, err
	// 	}
	//
	// 	s.Cache(acc)
	//
	// 	return acc, nil
	// }
}

// UnsafeLookupByName retrieves an account by its name. It's unsafe to use this method because
// it's reading from the Redis database.
func (s *AccountService) UnsafeLookupByName(name string) (*account.Account, error) {
	if acc := s.LookupByName(name); acc != nil {
		return acc, nil
	}

	val, err := common.RedisClient.Get(context.Background(), "primal%names:"+strings.ToLower(name)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, errors.New("key does not exists")
	}

	if err != nil {
		return nil, err
	}

	if val == "" {
		return nil, errors.New("empty value")
	}

	acc := &account.Account{}
	if acc.UnmarshalString(val) != nil {
		return nil, err
	}

	s.Cache(acc)

	return acc, nil
}

// UpdateName updates the name of an account.
func (s *AccountService) UpdateName(oldName, newName, id string) {
	s.accountsIdMu.Lock()

	delete(s.accountsId, strings.ToLower(oldName))
	s.accountsId[strings.ToLower(newName)] = id

	s.accountsIdMu.Unlock()
}

// Cache caches an account.
func (s *AccountService) Cache(a *account.Account) {
	s.accountsMu.Lock()
	s.accounts[a.Id()] = a
	s.accountsMu.Unlock()

	s.accountsIdMu.Lock()
	s.accountsId[strings.ToLower(a.Name())] = a.Id()
	s.accountsIdMu.Unlock()
}

// RedisCache caches an account into the Redis database.
func (s *AccountService) RedisCache(acc *account.Account) error {
	if common.RedisClient == nil {
		return errors.New("redis client not found")
	}

	pip := common.RedisClient.Pipeline()
	if pip == nil {
		return errors.New("redis pipeline not found")
	}

	pip.Set(context.Background(), "primal%ids:"+acc.Id(), acc.MarshalString(), 72*time.Hour)
	pip.Set(context.Background(), "primal%names:"+strings.ToLower(acc.Name()), acc.MarshalString(), 72*time.Hour)

	if _, err := pip.Exec(context.Background()); err != nil {
		return err
	}

	return nil
}

// Update updates an account.
func (s *AccountService) Update(acc *account.Account) error {
	if s.col == nil {
		return errors.New("service not hooked to the database")
	}

	res, err := s.col.UpdateOne(
		context.TODO(),
		bson.D{{"_id", acc.Id()}},
		bson.D{{"$set", acc.Marshal()}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return err
	}

	if res.UpsertedCount > 0 {
		common.Log.Printf("Account %s was inserted", acc.Id())
	} else {
		common.Log.Printf("Account %s was updated", acc.Id())
	}

	if err := s.RedisCache(acc); err != nil {
		return err
	}

	return nil
}

func (s *AccountService) lookupAtRedis(k, v string) (*account.Account, error) {
	val, err := common.RedisClient.Get(context.Background(), k+v).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else if val == "" {
		return nil, errors.New("empty value")
	} else {
		acc := &account.Account{}
		if acc.UnmarshalString(val) != nil {
			return nil, err
		}

		return acc, nil
	}
}

func (s *AccountService) lookupAtMongo(k, v string) (*account.Account, error) {
	var acc account.Account

	if err := s.col.FindOne(context.TODO(), bson.D{{k, v}}).Decode(&acc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &acc, nil
}

// Hook hooks the account service to the database.
func (s *AccountService) Hook(db *mongo.Database) error {
	if s.col != nil {
		return errors.New("service already hooked to the database")
	}

	s.col = db.Collection("accounts")

	return nil
}

// Account returns the account service.
func Account() *AccountService {
	return accountService
}

var accountService = &AccountService{
	accounts:   make(map[string]*account.Account),
	accountsId: make(map[string]string),
}
