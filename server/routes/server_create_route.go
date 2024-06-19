package routes

import (
	"context"
	"fmt"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/server"
	"github.com/holypvp/primal/server/pubsub"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

func ServerCreateRoute(c echo.Context) error {
	serverId := c.Param("id")
	if serverId == "" {
		return c.String(http.StatusBadRequest, "No ID found")
	}

	serverInfo := server.Service().LookupById(serverId)
	if serverInfo != nil {
		return c.String(http.StatusConflict, fmt.Sprintf("Server %s already exists", serverId))
	}

	port := c.Param("port")
	if port == "" {
		return c.String(http.StatusBadRequest, fmt.Sprintf("No port found for server %s", serverId))
	}

	portNum, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Invalid port for server %s", serverId))
	}

	if server.Service().LookupByPort(portNum) != nil {
		return c.String(http.StatusConflict, fmt.Sprintf("Port %d is already in use", portNum))
	}

	payload, err := common.WrapPayload("API_SERVER_CREATE", pubsub.NewServerCreatePacket(serverId, portNum))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to marshal packet")
	}

	serverInfo = server.NewServerInfo(serverId, portNum)
	serverInfo.SetInitialTime(time.Now().UnixMilli())

	server.Service().AppendServer(serverInfo)

	// Save the model into MongoDB but in a goroutine, so it doesn't block the main thread
	// Here you have the difference between the two snippets
	// Main thread ms = +133ms / Goroutine ms = 63ms average
	go server.SaveModel(serverInfo.ToModel())

	err = common.RedisClient.Publish(context.Background(), common.RedisChannel, payload).Err()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to publish packet")
	}

	return c.String(http.StatusOK, fmt.Sprintf("Server %s created on port %d", serverId, portNum))
}
