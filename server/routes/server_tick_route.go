package routes

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/server"
	"github.com/holypvp/primal/server/request"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ServerTickRoute(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return common.HTTPError(http.StatusBadRequest, "No server ID found")
	}

	serverInfo := server.Service().LookupById(id)
	if serverInfo == nil {
		return common.HTTPError(http.StatusNoContent, "Server not found")
	}

	body := &request.ServerTickBodyRequest{}
	err := json.NewDecoder(c.Request().Body).Decode(body)
	if err != nil {
		return common.HTTPError(http.StatusBadRequest, errors.Join(errors.New("failed to decode body"), err).Error())
	}

	serverInfo.SetPlayersCount(body.PlayersCount)
	serverInfo.SetHeartbeat(body.Heartbeat)
	serverInfo.SetPlayers(body.Players)

	serverInfo.SetActiveThreads(body.ActiveThreads)
	serverInfo.SetDaemonThreads(body.DaemonThreads)

	serverInfo.SetTicksPerSecond(body.TicksPerSecond)
	serverInfo.SetFullTicks(body.FullTicks)

	// TODO: This have performance issues because it's blocking the main thread
	// so I prefer make the wrapper and publish in a goroutine
	payload, err := common.WrapPayload("API_SERVER_TICK", body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to wrap payload").SetInternal(err)
	}

	err = common.RedisClient.Publish(context.Background(), common.RedisChannel, payload).Err()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to publish payload").SetInternal(err)
	}

	return c.String(http.StatusOK, "Server ticked")
}
