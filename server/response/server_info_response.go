package response

import (
	"github.com/holypvp/primal/server/model"
)

type ServerInfoResponse struct {
	Id     string   `json:"id"`
	Port   int64    `json:"port"`
	Groups []string `json:"bgroups"`

	PlayersCount   int      `json:"players-count"`
	MaxSlots       int64    `json:"max-slots"`
	Heartbeat      int64    `json:"heartbeat"`
	Players        []string `json:"players"`
	BungeeCord     bool     `json:"bungee-cord"`
	OnlineMode     bool     `json:"online-mode"`
	ActiveThreads  int      `json:"active-threads"`
	DaemonThreads  int      `json:"daemon-threads"`
	Motd           string   `json:"motd"`
	TicksPerSecond float64  `json:"ticks-per-second"`
	Directory      string   `json:"directory"`
	FullTicks      float64  `json:"full-ticks"`
	InitialTime    int64    `json:"initial-time"`
	Plugins        []string `json:"plugins"`
}

func NewServerInfoResponse(serverInfo *model.ServerInfo) ServerInfoResponse {
	return ServerInfoResponse{
		Id:             serverInfo.Id(),
		Port:           serverInfo.Port(),
		Groups:         serverInfo.Groups(),
		PlayersCount:   serverInfo.PlayersCount(),
		MaxSlots:       serverInfo.MaxSlots(),
		Heartbeat:      serverInfo.Heartbeat(),
		Players:        serverInfo.Players(),
		BungeeCord:     serverInfo.BungeeCord(),
		OnlineMode:     serverInfo.OnlineMode(),
		ActiveThreads:  serverInfo.ActiveThreads(),
		DaemonThreads:  serverInfo.DaemonThreads(),
		Motd:           serverInfo.Motd(),
		TicksPerSecond: serverInfo.TicksPerSecond(),
		Directory:      serverInfo.Directory(),
		FullTicks:      serverInfo.FullTicks(),
		InitialTime:    serverInfo.InitialTime(),
		Plugins:        serverInfo.Plugins(),
	}
}
