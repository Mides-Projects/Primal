package routes

import (
	"context"
	"encoding/json"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/server"
	"github.com/holypvp/primal/server/request"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ServerTickRoute(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.String(http.StatusBadRequest, "Server ID is required")
	}

	serverInfo := server.Service().LookupById(id)
	if serverInfo == nil {
		return c.String(http.StatusNoContent, "Server not found")
	}

	body := &request.ServerTickBody{}
	err := json.NewDecoder(c.Request().Body).Decode(body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}

	serverInfo.SetPlayersCount(body.PlayersCount)
	serverInfo.SetHeartbeat(body.Heartbeat)
	serverInfo.SetPlayers(body.Players)

	serverInfo.SetActiveThreads(body.ActiveThreads)
	serverInfo.SetDaemonThreads(body.DaemonThreads)

	serverInfo.SetTicksPerSecond(body.TicksPerSecond)
	serverInfo.SetFullTicks(body.FullTicks)

	payload, err := common.WrapPayload("API_SERVER_TICK", body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to wrap payload")
	}

	err = common.RedisClient.Publish(context.Background(), common.RedisChannel, payload).Err()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to publish payload")
	}

	return c.String(http.StatusOK, "Server ticked")
}
