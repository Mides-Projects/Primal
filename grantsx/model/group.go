package model

import (
	"errors"
	"github.com/google/uuid"
)

type Group struct {
	id string // ID of the group

	name        string // Name of the group
	displayName string // Display name of the group

	charColor string // The character color of the group
	siteColor string // Site color of the group

	prefix string // Prefix string
	suffix string // Suffix string

	weight int32 // Weight of the group
	hidden bool  // Hidden status of the group

	permissions []string // Permissions of the group
	inherits    []string // Inherited bgroups

	// Metadata usually contains the next steps to be taken by the group.
	// queue_eviction_time: managed by the queue service
	// discord_id: managed by the discord service
	// price: managed by the store service
	metadata map[string]interface{}
}

func EmptyGroup(name string) *Group {
	return &Group{
		id:   uuid.New().String(),
		name: name,
	}
}

// Id returns the ID of the group.
func (g *Group) Id() string {
	return g.id
}

// Name returns the name of the group.
func (g *Group) Name() string {
	return g.name
}

// SetName sets the name of the group.
func (g *Group) SetName(name string) {
	g.name = name
}

// DisplayName returns the display name of the group.
func (g *Group) DisplayName() string {
	return g.displayName
}

// SetDisplayName sets the display name of the group.
func (g *Group) SetDisplayName(displayName string) {
	g.displayName = displayName
}

// CharColor returns the character color of the group.
func (g *Group) CharColor() string {
	return g.charColor
}

// SetCharColor sets the character color of the group.
func (g *Group) SetCharColor(charColor string) {
	g.charColor = charColor
}

// SiteColor returns the site color of the group.
func (g *Group) SiteColor() string {
	return g.siteColor
}

// SetSiteColor sets the site color of the group.
func (g *Group) SetSiteColor(siteColor string) {
	g.siteColor = siteColor
}

// Prefix returns the prefix string of the group.
func (g *Group) Prefix() string {
	return g.prefix
}

// SetPrefix sets the prefix string of the group.
func (g *Group) SetPrefix(prefix string) {
	g.prefix = prefix
}

// Suffix returns the suffix string of the group.
func (g *Group) Suffix() string {
	return g.suffix
}

// SetSuffix sets the suffix string of the group.
func (g *Group) SetSuffix(suffix string) {
	g.suffix = suffix
}

// Weight returns the weight of the group.
func (g *Group) Weight() int32 {
	return g.weight
}

// SetWeight sets the weight of the group.
func (g *Group) SetWeight(weight int32) {
	g.weight = weight
}

// Hidden returns the hidden status of the group.
func (g *Group) Hidden() bool {
	return g.hidden
}

// SetHidden sets the hidden status of the group.
func (g *Group) SetHidden(hidden bool) {
	g.hidden = hidden
}

// Permissions returns the permissions of the group.
func (g *Group) Permissions() []string {
	return g.permissions
}

// SetPermissions sets the permissions of the group.
func (g *Group) SetPermissions(permissions []string) {
	g.permissions = permissions
}

// Inherits returns the inherited bgroups of the group.
func (g *Group) Inherits() []string {
	return g.inherits
}

// SetInherits sets the inherited bgroups of the group.
func (g *Group) SetInherits(inherits []string) {
	g.inherits = inherits
}

// Metadata returns the metadata of the group.
func (g *Group) Metadata() map[string]interface{} {
	return g.metadata
}

// SetMetadata sets the metadata of the group.
func (g *Group) SetMetadata(metadata map[string]interface{}) {
	g.metadata = metadata
}

// Marshal marshals a group to a map.
func (g *Group) Marshal() map[string]interface{} {
	return map[string]interface{}{
		"id":           g.id,
		"name":         g.name,
		"display_name": g.displayName,
		"char_color":   g.charColor,
		"site_color":   g.siteColor,
		"prefix":       g.prefix,
		"suffix":       g.suffix,
		"weight":       g.weight,
		"hidden":       g.hidden,
		"permissions":  g.permissions,
		"inherits":     g.inherits,
		"metadata":     g.metadata,
	}
}

// Unmarshal unmarshals a group from a map.
func (g *Group) Unmarshal(body map[string]interface{}) error {
	id, ok := body["id"].(string)
	if !ok {
		return errors.New("missing 'id' field")
	}
	g.id = id

	name, ok := body["name"].(string)
	if !ok {
		return errors.New("missing 'name' field")
	}
	g.name = name

	displayName, ok := body["display_name"].(string)
	if !ok {
		return errors.New("missing 'display_name' field")
	}
	g.displayName = displayName

	charColor, ok := body["char_color"].(string)
	if !ok {
		return errors.New("missing 'char_color' field")
	}
	g.charColor = charColor

	siteColor, ok := body["site_color"].(string)
	if !ok {
		return errors.New("missing 'site_color' field")
	}
	g.siteColor = siteColor

	prefix, ok := body["prefix"].(string)
	if !ok {
		return errors.New("missing 'prefix' field")
	}
	g.prefix = prefix

	suffix, ok := body["suffix"].(string)
	if !ok {
		return errors.New("missing 'suffix' field")
	}
	g.suffix = suffix

	weight, ok := body["weight"].(int32)
	if !ok {
		return errors.New("missing 'weight' field")
	}
	g.weight = weight

	hidden, ok := body["hidden"].(bool)
	if !ok {
		return errors.New("missing 'hidden' field")
	}
	g.hidden = hidden

	if perms, ok := body["permissions"].([]string); ok {
		g.permissions = perms
	}

	if inherits, ok := body["inherits"].([]string); !ok {
		g.inherits = inherits
	}

	if metadata, ok := body["metadata"].(map[string]interface{}); !ok {
		g.metadata = metadata
	}

	return nil
}
