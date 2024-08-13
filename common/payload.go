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

func HTTPError(code int, message string) error {
	result, err := json.Marshal(map[string]interface{}{"code": code, "message": message})
	if err != nil {
		return echo.NewHTTPError(code, `{"code":500,"message":"Internal Server Error"}`)
	}

	return echo.NewHTTPError(code, string(result))
}
