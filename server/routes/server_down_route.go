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
		return echo.NewHTTPError(http.StatusBadRequest, "No ID found")
	}

	serverInfo := server.Service().LookupById(serverId)
	if serverInfo == nil {
		return echo.NewHTTPError(http.StatusNoContent, fmt.Sprintf("Server %s not found", serverId))
	}

	payload, err := common.WrapPayload("API_SERVER_DOWN", pubsub.NewServerStatusPacket(serverInfo.Id()))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to marshal packet").SetInternal(err)
	}

	err = common.RedisClient.Publish(context.Background(), common.RedisChannel, payload).Err()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to publish packet").SetInternal(err)
	}

	return c.String(http.StatusOK, fmt.Sprintf("Server %s is now down", serverInfo.Id()))
}
