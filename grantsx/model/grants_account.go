package model

import (
	"github.com/holypvp/primal/source/model"
)

type GrantsAccount struct {
	account *model.Account // Account of the account

	activeGrants  []Grant // Active grants of the account
	expiredGrants []Grant // Expired grants of the account
}

// Account returns the account of the account.
func (ga *GrantsAccount) Account() *model.Account {
	return ga.account
}

// ActiveGrants returns the active grants of the account.
func (ga *GrantsAccount) ActiveGrants() []Grant {
	return ga.activeGrants
}

// ExpiredGrants returns the expired grants of the account.
func (ga *GrantsAccount) ExpiredGrants() []Grant {
	return ga.expiredGrants
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
