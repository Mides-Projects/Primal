package common

import "encoding/json"

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

func ErrorResponse(code int, message string) string {
	result, err := json.Marshal(map[string]interface{}{
		"code":    code,
		"message": message,
	})
	if err != nil {
		return `{"code":500,"message":"Internal Server Error"}`
	}

	return string(result)
}
