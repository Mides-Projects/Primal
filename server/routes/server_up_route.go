package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/server"
	"github.com/holypvp/primal/server/pubsub"
	"github.com/holypvp/primal/server/request"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func ServerUpRoute(c echo.Context) error {
	serverId := c.Param("id")
	if serverId == "" {
		return c.String(http.StatusBadRequest, "Server ID is required")
	}

	serverInfo := server.Service().LookupById(serverId)
	if serverInfo == nil {
		return c.String(http.StatusNoContent, "Server not found")
	}

	body := &request.ServerUpBody{}
	err := json.NewDecoder(c.Request().Body).Decode(body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}

	serverInfo.SetDirectory(body.Directory)
	serverInfo.SetMotd(body.Motd)

	serverInfo.SetBungeeCord(body.BungeeCord)
	serverInfo.SetOnlineMode(body.OnlineMode)

	serverInfo.SetMaxSlots(body.MaxSlots)
	serverInfo.SetPlugins(body.Plugins)

	initialTime := serverInfo.InitialTime()
	if initialTime == 0 {
		return c.String(http.StatusBadRequest, "Server has not been created yet")
	}

	now := time.Now().UnixMilli()
	serverInfo.SetInitialTime(now)

	// Save the server model in a goroutine to avoid blocking the main thread
	go server.SaveModel(serverInfo.ToModel())

	payload, err := common.WrapPayload("API_SERVER_UP", pubsub.NewServerStatusPacket(serverId))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to wrap payload")
	}

	err = common.RedisClient.Publish(context.Background(), common.RedisChannel, payload).Err()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to publish payload")
	}

	return c.String(http.StatusOK, fmt.Sprintf("Server %s is now back up. After %d ms", serverId, now-serverInfo.Heartbeat()))
}
