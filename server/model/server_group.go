package model

import "errors"

type ServerGroup struct {
	id                    string
	metadata              map[string]interface{}
	announcements         []string
	announcementsInterval int64

	fallbackServerId *string
}

func NewServerGroup(id string) *ServerGroup {
	return &ServerGroup{
		id:                    id,
		metadata:              map[string]interface{}{},
		announcements:         []string{},
		announcementsInterval: 0,
		fallbackServerId:      nil,
	}
}

func (g *ServerGroup) Id() string {
	return g.id
}

func (g *ServerGroup) Metadata() map[string]interface{} {
	return g.metadata
}

func (g *ServerGroup) SetMetadata(metadata map[string]interface{}) {
	g.metadata = metadata
}

func (g *ServerGroup) Announcements() []string {
	return g.announcements
}

func (g *ServerGroup) SetAnnouncements(announcements []string) {
	g.announcements = announcements
}

func (g *ServerGroup) AnnouncementsInterval() int64 {
	return g.announcementsInterval
}

func (g *ServerGroup) SetAnnouncementsInterval(interval int64) {
	g.announcementsInterval = interval
}

func (g *ServerGroup) FallbackServerId() *string {
	return g.fallbackServerId
}

func (g *ServerGroup) SetFallbackServerId(fallbackServerId *string) {
	g.fallbackServerId = fallbackServerId
}

func (g *ServerGroup) AddAnnouncement(announcement string) {
	g.announcements = append(g.announcements, announcement)
}

func (g *ServerGroup) RemoveAnnouncement(announcement string) {
	for i, a := range g.announcements {
		if a == announcement {
			g.announcements = append(g.announcements[:i], g.announcements[i+1:]...)
		}
	}
}

func (g *ServerGroup) Marshal() map[string]interface{} {
	return map[string]interface{}{
		"id":                     g.id,
		"metadata":               g.metadata,
		"announcements":          g.announcements,
		"announcements_interval": g.announcementsInterval,
		"fallback_server_id":     g.fallbackServerId,
	}
}

func (g *ServerGroup) Unmarshal(data map[string]interface{}) error {
	id, ok := data["id"].(string)
	if !ok {
		return errors.New("id is not a string")
	}

	g.id = id

	metadata, ok := data["metadata"].(map[string]interface{})
	if !ok {
		return errors.New("metadata is not a map")
	}
	g.metadata = metadata

	announcements, ok := data["announcements"].([]string)
	if !ok {
		return errors.New("announcements is not a string array")
	}
	g.announcements = announcements

	if announcementsInterval, ok := data["announcements_interval"].(int64); !ok {
		return errors.New("announcements_interval is not an int64")
	} else {
		g.announcementsInterval = announcementsInterval
	}

	fallbackServerId, ok := data["fallback_server_id"].(*string)
	if !ok {
		return errors.New("fallback_server_id is not a string")
	}
	g.fallbackServerId = fallbackServerId

	return nil
}
