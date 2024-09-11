package service

import (
    "context"
    "errors"
    "github.com/holypvp/primal/common"
    "github.com/holypvp/primal/source/model"
    "github.com/redis/go-redis/v9"
    "strings"
    "sync"
)

type AccountService struct {
    accountsMu sync.RWMutex
    accounts   map[string]*model.Account

    accountsIdMu sync.RWMutex
    accountsId   map[string]string
}

func (s *AccountService) LookupById(id string) *model.Account {
    s.accountsMu.RLock()
    defer s.accountsMu.RUnlock()

    return s.accounts[id]
}

func (s *AccountService) LookupIdByName(name string) (string, bool) {
    s.accountsIdMu.RLock()
    defer s.accountsIdMu.RUnlock()

    id, ok := s.accountsId[strings.ToLower(name)]
    return id, ok
}

func (s *AccountService) Fetch(id string) (*model.Account, error) {
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
        acc := &model.Account{}
        if err = acc.UnmarshalString(val); err != nil {
            return nil, err
        }

        s.Cache(acc)

        return acc, nil
    }
}

// FetchAccountId retrieves an account by its ID from our database.
func (s *AccountService) FetchAccountId(name string) string {
    if id, ok := s.LookupIdByName(name); ok {
        return id
    }

    val, err := common.RedisClient.Get(context.Background(), "primal%names:"+strings.ToLower(name)).Result()
    if errors.Is(err, redis.Nil) {
        common.Log.Fatalf("Failed to fetch account ID '%v': Key does not exists", name)

        return ""
    }

    if err != nil {
        common.Log.Fatalf("Failed to fetch account ID '%v': %v", name, err)

        return ""
    }

    if val == "" {
        common.Log.Fatalf("Failed to fetch account ID: %v", name)

        return ""
    }

    acc := &model.Account{}
    if acc.UnmarshalString(val) != nil {
        return ""
    }

    s.Cache(acc)

    return acc.Id()
}

// Cache caches an account.
func (s *AccountService) Cache(a *model.Account) {
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
    accounts:   make(map[string]*model.Account),
    accountsId: make(map[string]string),
}
