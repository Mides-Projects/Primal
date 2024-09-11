package routes

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/server/model"
	"github.com/holypvp/primal/server/pubsub"
	"github.com/holypvp/primal/service"
	"log"
	"net/http"
	"strconv"
	"time"
)

func ServerCreateRoute(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return common.HTTPError(http.StatusBadRequest, "No server ID found")
	}

	serverInfo := service.Server().LookupById(id)
	if serverInfo != nil {
		return common.HTTPError(http.StatusConflict, fmt.Sprintf("Server %s already exists", id))
	}

	port := c.Params("port")
	if port == "" {
		return common.HTTPError(http.StatusBadRequest, "No port found")
	}

	portNum, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		return common.HTTPError(http.StatusBadRequest, "Invalid port number")
	}

	if service.Server().LookupByPort(portNum) != nil {
		return common.HTTPError(http.StatusConflict, fmt.Sprintf("Port %d already in use", portNum))
	}

	serverInfo = model.NewServerInfo(id, portNum)
	serverInfo.SetInitialTime(time.Now().UnixMilli())

	service.Server().CacheServer(serverInfo)

	// Save the model into MongoDB but in a goroutine, so it doesn't block the main thread
	// Here you have the difference between the two snippets
	// Main thread ms = +133ms / Goroutine ms = 63ms average
	// but small issue is I can't get the error from the goroutine
	go func() {
		if err = service.SaveModel(serverInfo.Id(), serverInfo.Marshal()); err != nil {
			common.Log.Errorf("Failed to save server %s: %v", serverInfo.Id(), err)
		}

		payload, err := common.WrapPayload("API_SERVER_CREATE", pubsub.NewServerCreatePacket(id, portNum))
		if err != nil {
			log.Fatal("Failed to marshal packet: ", err)
		}

		if err = common.RedisClient.Publish(context.Background(), common.RedisChannel, payload).Err(); err != nil {
			log.Fatal("Failed to publish packet: ", err)
		}
	}()

	return c.Status(http.StatusOK).SendString(fmt.Sprintf("Server %s created on port %d", id, portNum))
}

func Hook(app *fiber.App) {
	g := app.Group("/servers")

	g.Post("/:id/create/:port", ServerCreateRoute)
	g.Get("/:id/lookup", LookupServers)
	g.Patch("/:id/down", ServerDownRoute)
	g.Patch("/:id/up", ServerUpRoute)
	g.Patch("/:id/tick", ServerTickRoute)
}
