package routes

import (
	"context"
	"fmt"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/server"
	"github.com/holypvp/primal/server/model"
	"github.com/holypvp/primal/server/pubsub"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"time"
)

func ServerCreateRoute(c echo.Context) error {
	serverId := c.Param("id")
	if serverId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "No ID found")
	}

	serverInfo := server.Service().LookupById(serverId)
	if serverInfo != nil {
		return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("Server %s already exists", serverId))
	}

	port := c.Param("port")
	if port == "" {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("No port found for server %s", serverId))
	}

	portNum, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid port for server %s", serverId))
	}

	if server.Service().LookupByPort(portNum) != nil {
		return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("Port %d is already in use", portNum))
	}

	serverInfo = model.NewServerInfo(serverId, portNum)
	serverInfo.SetInitialTime(time.Now().UnixMilli())

	server.Service().CacheServer(serverInfo)

	// Save the model into MongoDB but in a goroutine, so it doesn't block the main thread
	// Here you have the difference between the two snippets
	// Main thread ms = +133ms / Goroutine ms = 63ms average
	// but small issue is I can't get the error from the goroutine
	go func() {
		if err = server.SaveModel(serverInfo.Id(), serverInfo.Marshal()); err != nil {
			common.Log.Errorf("Failed to save server %s: %v", serverInfo.Id(), err)
		}

		payload, err := common.WrapPayload("API_SERVER_CREATE", pubsub.NewServerCreatePacket(serverId, portNum))
		if err != nil {
			log.Fatal("Failed to marshal packet: ", err)
		}

		if err = common.RedisClient.Publish(context.Background(), common.RedisChannel, payload).Err(); err != nil {
			log.Fatal("Failed to publish packet: ", err)
		}
	}()

	return c.String(http.StatusOK, fmt.Sprintf("Server %s created on port %d", serverId, portNum))
}
