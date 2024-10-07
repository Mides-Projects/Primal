package player

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type PlayerInfo struct {
	id string // ID of account

	name        string // The name of account
	lastName    string // The last name of account
	displayName string // The display name of account

	operator bool // Operator of account
	online   bool // Online status of account

	highestGroup  string // Highest group of account
	currentServer string // Current server of account

	lastJoin time.Time
}

// Id returns the ID of the account.
func (pi *PlayerInfo) Id() string {
	return pi.id
}

// Name returns the name of the account.
func (pi *PlayerInfo) Name() string {
	return pi.name
}

// SetName sets the name of the account.
func (pi *PlayerInfo) SetName(name string) {
	pi.name = name
}

// LastName returns the last name of the account.
func (pi *PlayerInfo) LastName() string {
	return pi.lastName
}

// SetLastName sets the last name of the account.
func (pi *PlayerInfo) SetLastName(lastName string) {
	pi.lastName = lastName
}

// DisplayName returns the display name of the account.
func (pi *PlayerInfo) DisplayName() string {
	return pi.displayName
}

// SetDisplayName sets the display name of the account.
func (pi *PlayerInfo) SetDisplayName(displayName string) {
	pi.displayName = displayName
}

// Operator returns the operator of the account.
func (pi *PlayerInfo) Operator() bool {
	return pi.operator
}

// SetOperator sets the operator of the account.
func (pi *PlayerInfo) SetOperator(operator bool) {
	pi.operator = operator
}

// Online returns the online status of the account.
func (pi *PlayerInfo) Online() bool {
	return pi.online
}

// SetOnline sets the online status of the account.
func (pi *PlayerInfo) SetOnline(online bool) {
	pi.online = online
}

// HighestGroup returns the highest group of the account.
func (pi *PlayerInfo) HighestGroup() string {
	return pi.highestGroup
}

// SetHighestGroup sets the highest group of the account.
func (pi *PlayerInfo) SetHighestGroup(group string) {
	pi.highestGroup = group
}

// CurrentServer returns the current server of the account.
func (pi *PlayerInfo) CurrentServer() string {
	return pi.currentServer
}

// SetCurrentServer sets the current server of the account.
func (pi *PlayerInfo) SetCurrentServer(server string) {
	pi.currentServer = server
}

// LastJoin returns the last join time of the account.
func (pi *PlayerInfo) LastJoin() time.Time {
	return pi.lastJoin
}

// SetLastJoin sets the last join time of the account.
func (pi *PlayerInfo) SetLastJoin(join time.Time) {
	pi.lastJoin = join
}

// Marshal returns the account as a map.
func (pi *PlayerInfo) Marshal() map[string]interface{} {
	return map[string]interface{}{
		"_id":       pi.id,
		"name":      pi.name,
		"last_name": pi.lastName,
		"operator":  pi.operator,
	}
}

// Unmarshal unmarshals the account from a map.
func (pi *PlayerInfo) Unmarshal(body map[string]interface{}) error {
	id, ok := body["_id"].(string)
	if !ok {
		return errors.New("id is not pi string")
	}
	pi.id = id

	name, ok := body["name"].(string)
	if !ok {
		return errors.New("name is not pi string")
	}
	pi.name = name

	lastName, ok := body["last_name"].(string)
	if !ok {
		return errors.New("last_name is not pi string")
	}
	pi.lastName = lastName

	operator, ok := body["operator"].(bool)
	if !ok {
		return errors.New("operator is not pi bool")
	}
	pi.operator = operator

	return nil
}

// UnmarshalString unmarshals the account from a string.
func (pi *PlayerInfo) UnmarshalString(result string) error {
	body := strings.Split(result, ":")
	if len(body) != 4 {
		return errors.New("body is more than 4 elements")
	}

	pi.id = body[0]
	pi.name = body[1]
	pi.lastName = body[2]
	pi.operator = body[3] == "true"
	pi.online = false

	return nil
}

// MarshalString marshals the account to a string.
func (pi *PlayerInfo) MarshalString() string {
	return pi.id + ":" + pi.name + ":" + pi.lastName + ":" + strconv.FormatBool(pi.operator)
}

// String returns the account as a string.
func (pi *PlayerInfo) String() string {
	return fmt.Sprintf(
		"PlayerInfo{id=%s, name=%s, last_name=%s, operator=%t, display_name=%s, highest_group=%s}",
		pi.id,
		pi.name,
		pi.lastName,
		pi.operator,
		pi.displayName,
		pi.highestGroup)
}

// Empty returns an empty account.
func Empty(id, name string) *PlayerInfo {
	return &PlayerInfo{
		id:   id,
		name: name,
	}
}
