package server

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/server/request"
	"github.com/holypvp/primal/service"
	"net/http"
)

func ServerTickRoute(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return common.HTTPError(c, http.StatusBadRequest, "No server ID found")
	}

	serverInfo := service.Server().LookupById(id)
	if serverInfo == nil {
		return common.HTTPError(c, http.StatusNoContent, "Server not found")
	}

	body := &request.ServerTickBodyRequest{}
	if err := c.Bind().Body(body); err != nil {
		return common.HTTPError(c, http.StatusBadRequest, errors.Join(errors.New("failed to decode body"), err).Error())
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
		return common.HTTPError(c, http.StatusInternalServerError, "Failed to wrap payload")
	}

	err = common.RedisClient.Publish(context.Background(), common.RedisChannel, payload).Err()
	if err != nil {
		return common.HTTPError(c, http.StatusInternalServerError, "Failed to publish payload")
	}

	return c.Status(http.StatusOK).SendString("Server ticked")
}
