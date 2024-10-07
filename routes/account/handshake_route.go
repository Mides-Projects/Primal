package account

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/model/player"
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
	}

	if body.ServerName == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'server' body field",
		})
	}

	if body.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'name' body field",
		})
	}

	// TODO: Fix this because I need to let know him if the player already exists or not
	// if not exists, create a new account

	acc := service.Player().LookupById(id)
	if acc == nil && body.JoinedBefore {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "PlayerInfo for " + body.Name + " not found",
		})
	}

	empty := acc == nil
	if acc == nil {
		acc = player.Empty(id, "")
	}

	if acc.Name() != body.Name {
		oldName := acc.Name()
		acc.SetLastName(oldName)
		acc.SetName(body.Name)

		if !empty {
			service.Player().UpdateName(oldName, body.Name, acc.Id())
		} else {
			service.Player().Cache(acc, true)
		}

		go func() {
			if err := service.Player().Update(acc); err != nil {
				common.Log.Fatalf("Failed to update account: %s", err)
			}
		}()
	}

	if !acc.Online() {
		// TODO: Publish to the redis channel because the player joined!
	} else if acc.CurrentServer() != body.ServerName {
		// TODO: Publish to the redis channel because the player switched servers!
	}

	// Mark the account as online and set the current server
	acc.SetCurrentServer(body.ServerName)
	acc.SetOnline(true)
	acc.SetLastJoin(time.Now())

	// Im idiot
	return c.Status(http.StatusOK).JSON(acc.Marshal())
}

// Hook registers the route to the app
func Hook(app *fiber.App) {
	g := app.Group("/v1/account")
	g.Put("/:id/handshake", HandshakeRoute)
	g.Patch("/:id/update", UpdateRoute)
	g.Get("/:id/lookup", LookupRoute)
	g.Patch("/:id/quit", QuitRoute)
}
