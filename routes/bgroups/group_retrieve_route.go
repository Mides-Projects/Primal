package bgroups

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/routes/grants"
	"github.com/holypvp/primal/service"
	"net/http"
)

// retrieve retrieves all groups.
// This route is used to retrieve all groups. It's a simple route that
// returns all groups from the service.Groups() service.
func retrieve(c fiber.Ctx) error {
	values := service.Groups().All()

	var body []fiber.Map
	for _, g := range values {
		body = append(body, g.Marshal())
	}

	if body == nil {
		return c.Status(http.StatusNoContent).SendString("{}")
	}

	return c.Status(http.StatusOK).JSON(body)
}

func Hook(app *fiber.App) {
	// r.RouteNotFound(func(c *fiber.Ctx) error {

	g := app.Group("/v1/groups")
	g.Post("/:name/create/", create)
	g.Get("/", retrieve)
	// g.RouteNotFound("/*", func(_ echo.Context) error {
	// 	return common.HTTPError(echo.ErrLocked.Code, "This route is not available")
	// })

	gg := app.Group("/v1/grantsx")
	gg.Get("/:value/lookup/:type", grants.LookupRoute)
	gg.Get("/:value/lookup", grants.LookupRoute)
	gg.Post("/:name/create", grants.GrantsCreateRoute)
	// gg.RouteNotFound("/*", func(_ echo.Context) error {
	// 	return common.HTTPError(echo.ErrLocked.Code, "This route is not available")
	// })
}
