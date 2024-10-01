package account

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/model"
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

	var body map[string]interface{}
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse body: " + err.Error(),
		})
	}

	serverName, ok := body["server"].(string)
	if !ok || serverName == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'server' body field",
		})
	}

	name, ok := body["name"].(string)
	if !ok || name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'name' body field",
		})
	}

	exists := c.Query("exists")
	if exists == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'exists' query parameter",
		})
	}

	// TODO: Fix this because I need to let know him if the player already exists or not
	// if not exists, create a new account

	acc := service.Account().LookupById(id)
	if acc == nil && exists == "true" {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Account for " + name + " not found",
		})
	}

	empty := acc == nil
	if acc == nil {
		acc = model.Empty(id, "")
	}

	if acc.Name() != name {
		oldName := acc.Name()
		acc.SetLastName(oldName)
		acc.SetName(name)

		if !empty {
			service.Account().UpdateName(oldName, name, acc.Id())
		} else {
			service.Account().Cache(acc)
		}

		go func() {
			if err := service.Account().Update(acc); err != nil {
				common.Log.Fatalf("Failed to update account: %s", err)
			}
		}()
	}

	if !acc.Online() {
		// TODO: Publish to the redis channel because the player joined!
	} else if acc.CurrentServer() != serverName {
		// TODO: Publish to the redis channel because the player switched servers!
	}

	// Mark the account as online and set the current server
	acc.SetCurrentServer(serverName)
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
