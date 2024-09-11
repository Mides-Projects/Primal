package model

import (
	"github.com/holypvp/primal/source/model"
	"sync"
)

type GrantsAccount struct {
	account *model.Account // Account of the account

	activeMu     sync.RWMutex // Mutex for active grants
	activeGrants []*Grant     // Active grants of the account

	expiredMu     sync.RWMutex // Mutex for expired grants
	expiredGrants []*Grant     // Expired grants of the account
}

func EmptyGrantsAccount(account *model.Account) *GrantsAccount {
	return &GrantsAccount{
		account:       account,
		activeGrants:  make([]*Grant, 0),
		expiredGrants: make([]*Grant, 0),
	}
}

// Account returns the account of the account.
func (ga *GrantsAccount) Account() *model.Account {
	return ga.account
}

// ActiveGrants returns the active grants of the account.
func (ga *GrantsAccount) ActiveGrants() []*Grant {
	return ga.activeGrants
}

// AddActiveGrant adds a grant to the account.
func (ga *GrantsAccount) AddActiveGrant(g *Grant) {
	ga.activeMu.Lock()
	defer ga.activeMu.Unlock()

	ga.activeGrants = append(ga.activeGrants, g)
}

// ExpiredGrants returns the expired grants of the account.
func (ga *GrantsAccount) ExpiredGrants() []*Grant {
	return ga.expiredGrants
}

// AddExpiredGrant adds an expired grant to the account.
func (ga *GrantsAccount) AddExpiredGrant(g *Grant) {
	ga.expiredMu.Lock()
	defer ga.expiredMu.Unlock()

	ga.expiredGrants = append(ga.expiredGrants, g)
}

func (ga *GrantsAccount) Unmarshal(body map[string]interface{}) error {
	return nil
}

// Marshal returns the grants account as a map.
func (ga *GrantsAccount) Marshal() map[string]interface{} {
	return map[string]interface{}{
		"account":        ga.account.Marshal(),
		"active_grants":  ga.activeGrants,
		"expired_grants": ga.expiredGrants,
	}
}
