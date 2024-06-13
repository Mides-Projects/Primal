package model

import "github.com/holypvp/primal/server"

type ServerInfoModel struct {
	Id     string   `bson:"_id"`
	Port   int64    `bson:"port"`
	Groups []string `bson:"groups"`

	MaxSlots    int   `bson:"max_slots"`
	Heartbeat   int64 `bson:"heartbeat"`
	BungeeCord  bool  `bson:"bungee_cord"`
	OnlineMode  bool  `bson:"online_mode"`
	InitialTime int64 `bson:"initial_time"`
}

func WrapServerInfo(serverInfo *server.ServerInfo) ServerInfoModel {
	return ServerInfoModel{
		Id:          serverInfo.Id(),
		Port:        serverInfo.Port(),
		Groups:      serverInfo.Groups(),
		MaxSlots:    serverInfo.MaxSlots(),
		Heartbeat:   serverInfo.Heartbeat(),
		BungeeCord:  serverInfo.BungeeCord(),
		OnlineMode:  serverInfo.OnlineMode(),
		InitialTime: serverInfo.InitialTime(),
	}
}
