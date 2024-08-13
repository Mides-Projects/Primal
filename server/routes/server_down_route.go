package routes

import (
	"context"
	"fmt"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/server"
	"github.com/holypvp/primal/server/pubsub"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ServerDownRoute(c echo.Context) error {
	serverId := c.Param("id")
	if serverId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, "No server ID found"))
	}

	i := server.Service().LookupById(serverId)
	if i == nil {
		return echo.NewHTTPError(http.StatusNoContent, common.ErrorResponse(http.StatusNoContent, fmt.Sprintf("Server %s not found", serverId)))
	}

	go func() {
		payload, err := common.WrapPayload("API_SERVER_DOWN", pubsub.NewServerStatusPacket(i.Id()))
		if err != nil {
			common.Log.Fatal("Failed to marshal packet: ", err)
		}

		err = common.RedisClient.Publish(context.Background(), common.RedisChannel, payload).Err()
		if err != nil {
			common.Log.Errorf("Failed to publish packet: %v", err)
		}
	}()

	return c.String(http.StatusOK, fmt.Sprintf("Server %s is now down", i.Id()))
}
