package grantsx

import (
	"sync"
)

type Tracker struct {
	grants []*Grant
	mu     sync.Mutex
}

func EmptyTracker() *Tracker {
	return &Tracker{
		grants: make([]*Grant, 0),
	}
}

// ActiveGrants returns the active grantsx of the account.
func (t *Tracker) ActiveGrants() []*Grant {
	t.mu.Lock()
	defer t.mu.Unlock()

	var activeGrants []*Grant
	for _, g := range t.grants {
		if g.Expired() {
			continue
		}

		activeGrants = append(activeGrants, g)
	}

	return activeGrants
}

// AddActiveGrant adds a grant to the account.
func (t *Tracker) AddActiveGrant(g *Grant) {
	t.mu.Lock()
	t.grants = append(t.grants, g)
	t.mu.Unlock()
}

// ExpiredGrants returns the expired grantsx of the account.
func (t *Tracker) ExpiredGrants() []*Grant {
	t.mu.Lock()
	defer t.mu.Unlock()

	var expiredGrants []*Grant
	for _, g := range t.grants {
		if !g.Expired() {
			continue
		}

		expiredGrants = append(expiredGrants, g)
	}

	return expiredGrants
}

// AddExpiredGrant adds an expired grant to the account.
func (t *Tracker) AddExpiredGrant(g *Grant) {
	t.mu.Lock()
	t.grants = append(t.grants, g)
	t.mu.Unlock()
}

// Marshal returns the grants as a map.
// if filter is "active", only active grants are returned.
// if filter is empty, all grants are returned.
func (t *Tracker) Marshal(filter string) []map[string]interface{} {
	var body []map[string]interface{}
	if filter == "" {
		t.mu.Lock()
		for _, g := range t.grants {
			body = append(body, g.Marshal())
		}
		t.mu.Unlock()

		return body
	}

	for _, g := range t.ActiveGrants() {
		body = append(body, g.Marshal())
	}

	return body
}
