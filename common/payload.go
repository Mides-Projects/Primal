package common

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
)

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

type HTTPErrorPayload struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func HTTPError(code int, message string) error {
	return echo.NewHTTPError(code, HTTPErrorPayload{
		Code:    code,
		Message: message,
	})
}

func WrapPayload(pid string, payload interface{}) ([]byte, error) {
	return json.Marshal(NewPayload(pid, 0, payload))
}
