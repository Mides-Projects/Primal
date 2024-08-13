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
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, "No server ID found"))
	}

	si := server.Service().LookupById(serverId)
	if si == nil {
		return echo.NewHTTPError(http.StatusNoContent, common.ErrorResponse(http.StatusNoContent, fmt.Sprintf("Server %s not found", serverId)))
	}

	body := &request.ServerUpBodyRequest{}
	err := json.NewDecoder(c.Request().Body).Decode(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, "Failed to decode body")).SetInternal(err)
	}

	si.SetDirectory(body.Directory)
	si.SetMotd(body.Motd)

	si.SetBungeeCord(body.BungeeCord)
	si.SetOnlineMode(body.OnlineMode)

	si.SetMaxSlots(body.MaxSlots)
	si.SetPlugins(body.Plugins)

	initialTime := si.InitialTime()
	if initialTime == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, "Server has not been initialized"))
	}

	now := time.Now().UnixMilli()
	si.SetInitialTime(now)

	// Save the server model in a goroutine to avoid blocking the main thread
	go func() {
		if err = server.SaveModel(si.Id(), si.Marshal()); err != nil {
			common.Log.Errorf("Failed to save server %s: %v", si.Id(), err)
		}
	}()

	payload, err := common.WrapPayload("API_SERVER_UP", pubsub.NewServerStatusPacket(serverId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to wrap payload").SetInternal(err)
	}

	// TODO: Maybe we need do the publish in a goroutine too
	err = common.RedisClient.Publish(context.Background(), common.RedisChannel, payload).Err()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to publish payload").SetInternal(err)
	}

	return c.String(http.StatusOK, fmt.Sprintf("Server %s is now back up. After %d ms", serverId, now-si.Heartbeat()))
}
