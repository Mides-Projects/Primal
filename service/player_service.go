package service

import (
	"context"
	"errors"
	quark "github.com/Mides-Projects/Quark"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/model/player"
	"github.com/holypvp/primal/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"sync"
	"time"
)

type PlayerService struct {
	col          *mongo.Collection
	redisService redis.Service

	playersMu sync.RWMutex
	players   map[string]*player.PlayerInfo

	ttlCache *quark.Quark[*player.PlayerInfo]

	playersNameMu sync.RWMutex
	playersName   map[string]string
}

// DoTTLTick does a tick on the TTL cache.
func (s *PlayerService) DoTTLTick() {
	s.ttlCache.DoTick()
}

// LookupById retrieves an account by its ID. It's safe to use this method because
// it's only reading the map.
func (s *PlayerService) LookupById(id string) *player.PlayerInfo {
	s.playersMu.RLock()
	defer s.playersMu.RUnlock()

	if acc, ok := s.players[id]; ok {
		return acc
	}

	if acc, ok := s.ttlCache.Get(id); ok {
		return acc
	}

	return nil
}

// LookupByName retrieves an account by its name. It's safe to use this method because
// it's only reading the map.
func (s *PlayerService) LookupByName(name string) *player.PlayerInfo {
	s.playersNameMu.RLock()
	defer s.playersNameMu.RUnlock()

	if id, ok := s.playersName[strings.ToLower(name)]; ok {
		return s.LookupById(id)
	}

	return nil
}

// UnsafeLookupById retrieves an account by its ID. It's unsafe to use this method because
// it's reading from the Redis database.
func (s *PlayerService) UnsafeLookupById(id string, keep bool) (*player.PlayerInfo, error) {
	if acc := s.LookupById(id); acc != nil {
		return acc, nil
	}

	// These values are only cached for 72 hours, after that they are removed from redis
	// but still available in our mongo database.
	// If the account was fetch from database, it will be cached into redis to prevent further database calls in the next 72 hours.
	var (
		body map[string]interface{}
		err  error
	)
	if body, err = s.redisService.LookupJSON("ids:" + id); err != nil {
		return nil, err
	} else if body != nil {
		return s.wrap(body, keep)
	} else if s.col == nil {
		return nil, errors.New("service not hooked to the database")
	} else if err = s.col.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&body); err != nil {
		return nil, err
	} else {
		return s.wrap(body, keep)
	}
}

// UnsafeLookupByName retrieves an account by its name. It's unsafe to use this method because
// it's reading from the Redis database.
func (s *PlayerService) UnsafeLookupByName(name string, keep bool) (*player.PlayerInfo, error) {
	if acc := s.LookupByName(name); acc != nil {
		return acc, nil
	}

	var (
		body map[string]interface{}
	)

	if id, err := s.redisService.LookupString("names:" + strings.ToLower(name)); err != nil {
		return nil, err
	} else if id != "" {
		return s.UnsafeLookupById(id, keep)
	} else if s.col == nil {
		return nil, errors.New("service not hooked to the database")
	} else if err = s.col.FindOne(context.TODO(), bson.D{{"name", name}}).Decode(&body); err != nil {
		return nil, err
	} else {
		return s.wrap(body, keep)
	}
}

// UpdateName updates the name of an account.
func (s *PlayerService) UpdateName(oldName, newName, id string) {
	s.playersNameMu.Lock()

	delete(s.playersName, strings.ToLower(oldName))
	s.playersName[strings.ToLower(newName)] = id

	s.playersNameMu.Unlock()
}

func (s *PlayerService) wrap(body map[string]interface{}, keep bool) (*player.PlayerInfo, error) {
	acc := &player.PlayerInfo{}
	if err := acc.Unmarshal(body); err != nil {
		return nil, err
	}

	s.Cache(acc, keep)

	return acc, nil
}

// Cache caches an account.
func (s *PlayerService) Cache(pi *player.PlayerInfo, keep bool) {
	if s.ttlCache == nil {
		panic("service not hooked to the ttlCache")
	}

	if keep {
		s.playersMu.Lock()
		s.players[pi.Id()] = pi
		s.playersMu.Unlock()
	} else {
		s.ttlCache.Set(pi.Id(), pi)
	}

	s.playersNameMu.Lock()
	s.playersName[strings.ToLower(pi.Name())] = pi.Id()
	s.playersNameMu.Unlock()
}

// Invalidate invalidates an account.
func (s *PlayerService) Invalidate(acc *player.PlayerInfo) {
	s.ttlCache.Invalidate(acc.Id())

	s.playersMu.Lock()
	delete(s.players, acc.Id())
	s.playersMu.Unlock()

	s.playersNameMu.Lock()
	delete(s.playersName, strings.ToLower(acc.Name()))
	s.playersNameMu.Unlock()
}

// InvalidateTTL invalidates an account by its ID.
func (s *PlayerService) InvalidateTTL(id string) {
	s.ttlCache.Invalidate(id)
}

// Update updates an account.
func (s *PlayerService) Update(acc *player.PlayerInfo) error {
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
		common.Log.Printf("PlayerInfo %s was inserted", acc.Id())
	} else {
		common.Log.Printf("PlayerInfo %s was updated", acc.Id())
	}

	if err = s.redisService.StoreJSON("ids:"+acc.Id(), acc.Marshal(), 50*time.Hour); err != nil {
		return err
	} else if err = s.redisService.StoreString("names:"+strings.ToLower(acc.Name()), acc.Id(), 50*time.Hour); err != nil {
		return err
	}

	return nil
}

// Hook hooks the account service to the database.
func (s *PlayerService) Hook(db *mongo.Database) error {
	if s.col != nil {
		return errors.New("service already hooked to the database")
	}

	s.col = db.Collection("trackers")

	s.ttlCache = quark.New[*player.PlayerInfo](2*time.Hour, 2*time.Hour)
	s.ttlCache.SetListener(func(_ string, value *player.PlayerInfo, reason quark.Reason) {
		if reason == quark.ManualReason {
			return
		}

		s.playersNameMu.Lock()
		delete(s.playersName, strings.ToLower(value.Name()))
		s.playersNameMu.Unlock()
	})

	return nil
}

// Player returns the account service.
func Player() *PlayerService {
	return playerService
}

var playerService = &PlayerService{
	players:      make(map[string]*player.PlayerInfo),
	playersName:  make(map[string]string),
	redisService: redis.NewService("primal%"),
}
