package account

import (
	"errors"
	"strconv"
	"strings"
)

type Account struct {
	id string // ID of account

	name     string // The name of account
	lastName string // The last name of account

	operator bool // Operator of account
	online   bool // Online status of account

	currentServer string // Current server of account
}

// Id returns the ID of the account.
func (a *Account) Id() string {
	return a.id
}

// Name returns the name of the account.
func (a *Account) Name() string {
	return a.name
}

// SetName sets the name of the account.
func (a *Account) SetName(name string) {
	a.name = name
}

// LastName returns the last name of the account.
func (a *Account) LastName() string {
	return a.lastName
}

// SetLastName sets the last name of the account.
func (a *Account) SetLastName(lastName string) {
	a.lastName = lastName
}

// Operator returns the operator of the account.
func (a *Account) Operator() bool {
	return a.operator
}

// SetOperator sets the operator of the account.
func (a *Account) SetOperator(operator bool) {
	a.operator = operator
}

// Online returns the online status of the account.
func (a *Account) Online() bool {
	return a.online
}

// SetOnline sets the online status of the account.
func (a *Account) SetOnline(online bool) {
	a.online = online
}

// CurrentServer returns the current server of the account.
func (a *Account) CurrentServer() string {
	return a.currentServer
}

// SetCurrentServer sets the current server of the account.
func (a *Account) SetCurrentServer(server string) {
	a.currentServer = server
}

// Marshal returns the account as a map.
func (a *Account) Marshal() map[string]interface{} {
	return map[string]interface{}{
		"id":        a.id,
		"name":      a.name,
		"last_name": a.lastName,
		"operator":  a.operator,
	}
}

// Unmarshal unmarshals the account from a map.
func (a *Account) Unmarshal(body map[string]interface{}) error {
	id, ok := body["id"].(string)
	if !ok {
		return errors.New("id is not a string")
	}
	a.id = id

	name, ok := body["name"].(string)
	if !ok {
		return errors.New("name is not a string")
	}
	a.name = name

	lastName, ok := body["last_name"].(string)
	if !ok {
		return errors.New("last_name is not a string")
	}
	a.lastName = lastName

	operator, ok := body["operator"].(bool)
	if !ok {
		return errors.New("operator is not a bool")
	}
	a.operator = operator

	return nil
}

// UnmarshalString unmarshals the account from a string.
func (a *Account) UnmarshalString(result string) error {
	body := strings.Split(result, ":")
	if len(body) != 4 {
		return errors.New("body is more than 4 elements")
	}

	a.id = body[0]
	a.name = body[1]
	a.lastName = body[2]
	a.operator = body[3] == "true"
	a.online = false

	return nil
}

// MarshalString marshals the account to a string.
func (a *Account) MarshalString() string {
	return a.id + ":" + a.name + ":" + a.lastName + ":" + strconv.FormatBool(a.operator)
}

// Empty returns an empty account.
func Empty(id, name string) *Account {
	return &Account{
		id:   id,
		name: name,
	}
}
