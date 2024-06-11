package response

type ServerTickRequest struct {
	PlayersCount int      `json:"players-count"`
	MaxSlots     int      `json:"max-slots"`
	Players      []string `json:"players"`
	Groups       []string `json:"groups"`
}
