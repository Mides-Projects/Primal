package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/account"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/service"
	"net/http"
)

func AccountJoinRoute(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'id' parameter",
			"code":    http.StatusBadRequest,
		})
	}

	serverName := c.Params("server")
	if serverName == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'server' parameter",
			"code":    http.StatusBadRequest,
		})
	}

	name := c.Query("name")
	if name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'name' query parameter",
			"code":    http.StatusBadRequest,
		})
	}

	acc, err := service.Account().UnsafeLookupById(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to lookup account: " + err.Error(),
			"code":    http.StatusInternalServerError,
		})
	}

	if acc == nil {
		acc = account.Empty(id, "")
	}

	if acc.Name() != name {
		if acc.Name() != "" {
			service.Account().UpdateName(acc.Name(), name, acc.Id())
		} else {
			service.Account().Cache(acc)
		}

		acc.SetLastName(acc.Name())
		acc.SetName(name)

		go func() {
			if err = service.Account().Update(acc); err != nil {
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

	return c.Status(http.StatusOK).JSON(acc)
}

// Hook registers the route to the app
func Hook(app *fiber.App) {
	g := app.Group("/v1/account")
	g.Get("/:id/join/:name", AccountJoinRoute)
	g.Patch("/:id/update", AccountUpdateRoute)
	g.Delete("/:id/quit", AccountQuitRoute)
}
