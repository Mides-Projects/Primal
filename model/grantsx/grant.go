package grantsx

import "errors"

type Grant struct {
	id string // ID of grant

	identifier GrantIdentifier // Identifier of the grant

	addedBy string // ID of the player who added the grant
	addedAt string // Date when the grant was added

	expiresAt string // Date when the grant expires

	revokedBy string // ID of the player who revoked the grant
	revokedAt string // Date when the grant was revoked

	reason string // Reason for the grant

	scopes []string // Scopes of the grant
}

// Id returns the ID of the grant.
func (g *Grant) Id() string {
	return g.id
}

// Identifier returns the value of the grant.
func (g *Grant) Identifier() GrantIdentifier {
	return g.identifier
}

// AddedBy returns the ID of the player who added the grant.
func (g *Grant) AddedBy() string {
	return g.addedBy
}

// AddedAt returns the date when the grant was added.
func (g *Grant) AddedAt() string {
	return g.addedAt
}

// ExpiresAt returns the date when the grant expires.
func (g *Grant) ExpiresAt() string {
	return g.expiresAt
}

func (g *Grant) Expired() bool {
	if g.expiresAt == "" {
		return false
	}

	return false
}

// RevokedBy returns the ID of the player who revoked the grant.
func (g *Grant) RevokedBy() string {
	return g.revokedBy
}

// RevokedAt returns the date when the grant was revoked.
func (g *Grant) RevokedAt() string {
	return g.revokedAt
}

// Reason returns the reason for the grant.
func (g *Grant) Reason() string {
	return g.reason
}

// Scopes returns the scopes of the grant.
func (g *Grant) Scopes() []string {
	return g.scopes
}

func (g *Grant) Marshal() map[string]interface{} {
	return map[string]interface{}{
		"_id":        g.id,
		"identifier": g.identifier.Marshal(),
		"added_by":   g.addedBy,
		"added_at":   g.addedAt,
		"expires_at": g.expiresAt,
		"revoked_by": g.revokedBy,
		"revoked_at": g.revokedAt,
		"reason":     g.reason,
		"scopes":     g.scopes,
	}
}

// Unmarshal unmarshals a grant from a map.
func (g *Grant) Unmarshal(body map[string]interface{}) error {
	id, ok := body["_id"].(string)
	if !ok {
		return errors.New("_id is not a string")
	}
	g.id = id

	if idBody, ok := body["identifier"].(map[string]interface{}); !ok {
		return errors.New("identifier is not a map")
	} else if identifier, err := UnmarshalIdentifier(idBody); err != nil {
		return err
	} else {
		g.identifier = identifier
	}

	addedBy, ok := body["added_by"].(string)
	if !ok {
		return errors.New("added_by is not a string")
	}
	g.addedBy = addedBy

	addedAt, ok := body["added_at"].(string)
	if !ok {
		return errors.New("added_at is not a string")
	}
	g.addedAt = addedAt

	expiresAt, ok := body["expires_at"].(string)
	if !ok {
		return errors.New("expires_at is not a string")
	}
	g.expiresAt = expiresAt

	revokedBy, ok := body["revoked_by"].(string)
	if !ok {
		return errors.New("revoked_by is not a string")
	}
	g.revokedBy = revokedBy

	revokedAt, ok := body["revoked_at"].(string)
	if !ok {
		return errors.New("revoked_at is not a string")
	}
	g.revokedAt = revokedAt

	reason, ok := body["reason"].(string)
	if !ok {
		return errors.New("reason is not a string")
	}
	g.reason = reason

	if scopes, ok := body["scopes"].([]string); ok {
		g.scopes = scopes
	}

	return nil
}
