package account

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/model/player"
	"github.com/holypvp/primal/protocol"
	"github.com/holypvp/primal/redis"
	"github.com/holypvp/primal/service"
	"net/http"
	"time"
)

func HandshakeRoute(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'id' parameter",
		})
	}

	var body player.HandshakeBodyRequest
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse body: " + err.Error(),
		})
	} else if body.ServerName == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'server' body field",
		})
	} else if body.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'name' body field",
		})
	}

	pi := service.Player().LookupById(id)
	if pi == nil && body.JoinedBefore {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "PlayerInfo for " + body.Name + " not found",
		})
	}

	if pi == nil {
		pi = player.Empty(id, "")
	} else {
		service.Player().InvalidateTTL(pi.Id())
	}

	oldName := pi.Name()
	if oldName != body.Name {
		pi.SetLastName(oldName)
		pi.SetName(body.Name)

		if oldName != "" {
			service.Player().UpdateName(oldName, body.Name, id)
		} else {
			service.Player().Cache(pi, true)
		}

		go func() {
			if err := service.Player().Update(pi); err != nil {
				common.Log.Fatalf("Failed to update account: %s", err)
			}
		}()
	}

	var packet protocol.Packet
	if !pi.Online() {
		packet = &protocol.PlayerJoinedNetwork{
			Username:   pi.Name(),
			XUID:       pi.Id(),
			ServerName: body.ServerName,
		}
	} else if pi.CurrentServer() != body.ServerName {
		packet = &protocol.PlayerChangedServer{
			Username:      pi.Name(),
			XUID:          pi.Id(),
			OldServerName: pi.CurrentServer(),
			NewServerName: body.ServerName,
		}
	}

	if packet != nil {
		go redis.Publish(packet)
	}

	// Mark the account as online and set the current server
	pi.SetCurrentServer(body.ServerName)
	pi.SetOnline(true)
	pi.SetLastJoin(time.Now())

	// Im idiot
	return c.Status(http.StatusOK).JSON(pi.Marshal())
}

// Hook registers the route to the app
func Hook(app *fiber.App) {
	g := app.Group("/v1/account")
	g.Put("/:id/handshake", HandshakeRoute)
	g.Patch("/:id/update", UpdateRoute)
	g.Get("/:id/lookup", LookupRoute)
	g.Patch("/:id/quit", QuitRoute)
}
