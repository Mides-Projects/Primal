package response

type ServerTickRequest struct {
	Port           int64    `json:"port"`
	Groups         []string `json:"groups"`
	PlayersCount   int      `json:"players-count"`
	MaxSlots       int      `json:"max-slots"`
	Heartbeat      int64    `json:"heartbeat"`
	Players        []string `json:"players"`
	BungeeCord     bool     `json:"bungeecord"`
	OnlineMode     bool     `json:"online-mode"`
	ActiveThreads  int      `json:"active-threads"`
	DaemonThreads  int      `json:"daemon-threads"`
	Motd           *string  `json:"motd"`
	TicksPerSecond float64  `json:"ticks-per-second"`
	Directory      string   `json:"directory"`
	FullTicks      int64    `json:"full-ticks"`
	InitialTime    int64    `json:"initial-time"`
	Plugins        []string `json:"plugins"`
}
