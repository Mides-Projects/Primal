package grantsx

import (
	"github.com/holypvp/primal/model"
	"sync"
)

type Tracker struct {
	account *model.Account // Account of the account

	activeMu     sync.RWMutex // Mutex for active grantsx
	activeGrants []*Grant     // Active grantsx of the account

	expiredMu     sync.RWMutex // Mutex for expired grantsx
	expiredGrants []*Grant     // Expired grantsx of the account
}

func EmptyGrantsAccount(account *model.Account) *Tracker {
	return &Tracker{
		account:       account,
		activeGrants:  make([]*Grant, 0),
		expiredGrants: make([]*Grant, 0),
	}
}

// Account returns the account of the account.
func (ga *Tracker) Account() *model.Account {
	return ga.account
}

// ActiveGrants returns the active grantsx of the account.
func (ga *Tracker) ActiveGrants() []*Grant {
	return ga.activeGrants
}

// AddActiveGrant adds a grant to the account.
func (ga *Tracker) AddActiveGrant(g *Grant) {
	ga.activeMu.Lock()
	defer ga.activeMu.Unlock()

	ga.activeGrants = append(ga.activeGrants, g)
}

// ExpiredGrants returns the expired grantsx of the account.
func (ga *Tracker) ExpiredGrants() []*Grant {
	return ga.expiredGrants
}

// AddExpiredGrant adds an expired grant to the account.
func (ga *Tracker) AddExpiredGrant(g *Grant) {
	ga.expiredMu.Lock()
	defer ga.expiredMu.Unlock()

	ga.expiredGrants = append(ga.expiredGrants, g)
}
