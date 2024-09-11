package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/common"
	"net/http"
)

func accountLoadRoute(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return common.HTTPError(c, http.StatusBadRequest, "No account ID found")
	}

	// TODO: First check if the account is cached in the RAM, if it is, call a Handler to let know the api
	// the account changed of server, so we can send a broadcast to network.
	// acc, err := service.Account().UnsafeLookupById(id)

	// Send a json with the account info
	return c.Status(http.StatusOK).SendString("Account loaded")
}

func Hook(app *fiber.App) {
	g := app.Group("/v1/account")
	g.Get("/:id/load", accountLoadRoute)
}
