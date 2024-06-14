package request

type ServerTickBody struct {
	PlayersCount   int      `json:"players-count"`
	Heartbeat      int64    `json:"heartbeat"`
	Players        []string `json:"players"`
	ActiveThreads  int      `json:"active-threads"`
	DaemonThreads  int      `json:"daemon-threads"`
	TicksPerSecond float64  `json:"ticks-per-second"`
	FullTicks      float64  `json:"full-ticks"`
}
