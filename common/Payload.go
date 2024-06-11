package common

type Payload struct {
	PID     string      `json:"pid"`
	From    int64       `json:"from"`
	Payload interface{} `json:"payload"`
}

func NewPayload(pid string, from int64, payload interface{}) Payload {
	return Payload{
		PID:     pid,
		From:    from,
		Payload: payload,
	}
}
