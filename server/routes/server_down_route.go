package routes

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/server/pubsub"
	"github.com/holypvp/primal/service"
	"net/http"
)

func ServerDownRoute(c fiber.Ctx) error {
	serverId := c.Params("id")
	if serverId == "" {
		return common.HTTPError(c, http.StatusBadRequest, "No server ID found")
	}

	i := service.Server().LookupById(serverId)
	if i == nil {
		return common.HTTPError(c, http.StatusNoContent, fmt.Sprintf("Server %s not found", serverId))
	}

	go func() {
		payload, err := common.WrapPayload("API_SERVER_DOWN", pubsub.NewServerStatusPacket(i.Id()))
		if err != nil {
			common.Log.Fatal("Failed to marshal packet: ", err)
		}

		err = common.RedisClient.Publish(context.Background(), common.RedisChannel, payload).Err()
		if err != nil {
			common.Log.Fatalf("Failed to publish packet: %v", err)
		}
	}()

	return c.Status(http.StatusOK).SendString(fmt.Sprintf("Server %s is now down", i.Id()))
}
