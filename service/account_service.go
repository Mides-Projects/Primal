package service

import (
	"context"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/source/model"
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

// FetchAccountId retrieves an account by its ID from our database.
func (s *AccountService) FetchAccountId(name string) string {
	s.accountsIdMu.RLock()
	defer s.accountsIdMu.RUnlock()

	if id, ok := s.accountsId[strings.ToLower(name)]; ok {
		return id
	}

	res := common.RedisClient.Get(context.Background(), "primal%sources:"+strings.ToLower(name)).String()
	if res == "" {
		common.Log.Errorf("Failed to fetch account ID: %v", name)

		return ""
	}

	acc := &model.Account{}
	if acc.UnmarshalString(res) != nil {
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
}

// Account returns the account service.
func Account() *AccountService {
	return accountService
}

var accountService = &AccountService{
	accounts:   make(map[string]*model.Account),
	accountsId: make(map[string]string),
}
