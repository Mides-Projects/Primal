package player

type HandshakeBodyRequest struct {
    JoinedBefore bool   `json:"joined_before"`
    ServerName   string `json:"server"`
    Name         string `json:"name"`
}

type UpdateBodyRequest struct {
    HighestGroup string `json:"highest_group"`
    DisplayName  string `json:"display_name"`
    Timestamp    int64  `json:"timestamp"`
    Operator     bool   `json:"operator"`
    ServerName   string `json:"server"`
    Name         string `json:"name"`
}
