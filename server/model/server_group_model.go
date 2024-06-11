package model

type ServerGroupModel struct {
	Id                    string                 `bson:"_id"`
	Metadata              map[string]interface{} `bson:"metadata"`
	Announcements         []string               `bson:"announcements"`
	AnnouncementsInterval int64                  `bson:"announcements_interval"`
	FallbackServerId      *string                `bson:"fallback_server_id"`
}
