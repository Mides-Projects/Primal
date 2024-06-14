package model

type ServerInfoModel struct {
	Id     string   `bson:"_id"`
	Port   int64    `bson:"port"`
	Groups []string `bson:"groups"`

	MaxSlots    int64 `bson:"max_slots"`
	Heartbeat   int64 `bson:"heartbeat"`
	BungeeCord  bool  `bson:"bungee_cord"`
	OnlineMode  bool  `bson:"online_mode"`
	InitialTime int64 `bson:"initial_time"`
}
