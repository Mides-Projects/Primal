package grantsx

import (
    "sync"
)

type Tracker struct {
    activeMu     sync.RWMutex // Mutex for active grantsx
    activeGrants []*Grant     // Active grantsx of the account

    expiredMu     sync.RWMutex // Mutex for expired grantsx
    expiredGrants []*Grant     // Expired grantsx of the account
}

func EmptyTracker() *Tracker {
    return &Tracker{
        activeGrants:  make([]*Grant, 0),
        expiredGrants: make([]*Grant, 0),
    }
}

// ActiveGrants returns the active grantsx of the account.
func (t *Tracker) ActiveGrants() []*Grant {
    return t.activeGrants
}

// AddActiveGrant adds a grant to the account.
func (t *Tracker) AddActiveGrant(g *Grant) {
    t.activeMu.Lock()
    defer t.activeMu.Unlock()

    t.activeGrants = append(t.activeGrants, g)
}

// ExpiredGrants returns the expired grantsx of the account.
func (t *Tracker) ExpiredGrants() []*Grant {
    return t.expiredGrants
}

// AddExpiredGrant adds an expired grant to the account.
func (t *Tracker) AddExpiredGrant(g *Grant) {
    t.expiredMu.Lock()
    defer t.expiredMu.Unlock()

    t.expiredGrants = append(t.expiredGrants, g)
}

// Marshal returns the grants as a map.
// if filter is "active", only active grants are returned.
// if filter is empty, all grants are returned.
func (t *Tracker) Marshal(filter string) []map[string]interface{} {
    var body []map[string]interface{}
    for _, g := range t.activeGrants {
        body = append(body, g.Marshal())
    }

    if filter == "active" {
        return body
    }

    for _, g := range t.expiredGrants {
        body = append(body, g.Marshal())
    }

    return body
}
