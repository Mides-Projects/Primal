package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/server/pubsub"
	"github.com/holypvp/primal/server/request"
	"github.com/holypvp/primal/service"
	"net/http"
	"time"
)

func ServerUpRoute(c fiber.Ctx) error {
	serverId := c.Params("id")
	if serverId == "" {
		return common.HTTPError(c, http.StatusBadRequest, "No server ID found")
	}

	si := service.Server().LookupById(serverId)
	if si == nil {
		return common.HTTPError(c, http.StatusNoContent, fmt.Sprintf("Server %s not found", serverId))
	}

	body := &request.ServerUpBodyRequest{}
	if err := c.Bind().Body(body); err != nil {
		return common.HTTPError(c, http.StatusBadRequest, errors.Join(errors.New("failed to decode body"), err).Error())
	}

	si.SetDirectory(body.Directory)
	si.SetMotd(body.Motd)

	si.SetBungeeCord(body.BungeeCord)
	si.SetOnlineMode(body.OnlineMode)

	si.SetMaxSlots(body.MaxSlots)
	si.SetPlugins(body.Plugins)

	initialTime := si.InitialTime()
	if initialTime == 0 {
		return common.HTTPError(c, http.StatusBadRequest, "Server has not been initialized")
	}

	now := time.Now().UnixMilli()
	si.SetInitialTime(now)

	// Save the server model in a goroutine to avoid blocking the main thread
	go func() {
		if err := service.SaveModel(si.Id(), si.Marshal()); err != nil {
			common.Log.Fatalf("Failed to save server %s: %v", si.Id(), err)
		}
	}()

	payload, err := common.WrapPayload("API_SERVER_UP", pubsub.NewServerStatusPacket(serverId))
	if err != nil {
		return common.HTTPError(c, http.StatusInternalServerError, "Failed to wrap payload")
	}

	// TODO: Maybe we need do the publish in a goroutine too
	err = common.RedisClient.Publish(context.Background(), common.RedisChannel, payload).Err()
	if err != nil {
		return common.HTTPError(c, http.StatusInternalServerError, "Failed to publish payload")
	}

	return c.Status(http.StatusOK).SendString(fmt.Sprintf("Server %s is now back up. After %d ms", serverId, now-si.Heartbeat()))
}
