package service

import (
	"context"
	"errors"
	"github.com/holypvp/primal/account"
	"github.com/holypvp/primal/common"
	"github.com/redis/go-redis/v9"
	"strings"
	"sync"
)

type AccountService struct {
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

	val, err := common.RedisClient.Get(context.Background(), "primal%ids:"+id).Result()
	if errors.Is(err, redis.Nil) {
		return nil, errors.New("key does not exists")
	} else if err != nil {
		return nil, err
	} else if val == "" {
		return nil, errors.New("empty value")
	} else {
		acc := &account.Account{}
		if err = acc.UnmarshalString(val); err != nil {
			return nil, err
		}

		s.Cache(acc)

		return acc, nil
	}
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

// Cache caches an account.
func (s *AccountService) Cache(a *account.Account) {
	s.accountsMu.Lock()
	s.accounts[a.Id()] = a
	s.accountsMu.Unlock()

	s.accountsIdMu.Lock()
	s.accountsId[strings.ToLower(a.Name())] = a.Id()
	s.accountsIdMu.Unlock()
}

// Account returns the account service.
func Account() *AccountService {
	return accountService
}

var accountService = &AccountService{
	accounts:   make(map[string]*account.Account),
	accountsId: make(map[string]string),
}
