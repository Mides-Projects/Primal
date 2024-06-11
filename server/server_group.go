package server

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
