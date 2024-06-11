package response

type ServerInfoResponse struct {
	Id     string   `json:"id"`
	Port   int64    `json:"port"`
	Groups []string `json:"groups"`

	PlayersCount   int      `json:"players_count"`
	MaxSlots       int      `json:"max_slots"`
	Heartbeat      int64    `json:"heartbeat"`
	Players        []string `json:"players"`
	BungeeCord     bool     `json:"bungee_cord"`
	OnlineMode     bool     `json:"online_mode"`
	ActiveThreads  int      `json:"active_threads"`
	DaemonThreads  int      `json:"daemon_threads"`
	Motd           *string  `json:"motd"`
	TicksPerSecond float64  `json:"ticks_per_second"`
	Directory      string   `json:"directory"`
	FullTicks      int64    `json:"full_ticks"`
	InitialTime    int64    `json:"initial_time"`
	Plugins        []string `json:"plugins"`
}
