package response

type ServerGroupResponse struct {
	Id                    string                 `json:"_id"`
	Metadata              map[string]interface{} `json:"metadata"`
	Announcements         []string               `json:"announcements"`
	AnnouncementsInterval int64                  `json:"announcements_interval"`
	FallbackServerId      *string                `json:"fallback_server_id"`
}
