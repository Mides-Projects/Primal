package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/service"
	"net/http"
)

func groupsRetrieveRoute(c fiber.Ctx) error {
	values := service.Groups().All()

	body := make([]map[string]interface{}, len(values))
	for _, g := range values {
		body = append(body, g.Marshal())
	}

	return c.Status(http.StatusOK).JSON(body)
}

func Hook(app *fiber.App) {
	// r.RouteNotFound(func(c *fiber.Ctx) error {

	g := app.Group("/v1/groups")
	g.Post("/:name/create/", groupCreateRoute)
	g.Get("/", groupsRetrieveRoute)
	// g.RouteNotFound("/*", func(_ echo.Context) error {
	// 	return common.HTTPError(echo.ErrLocked.Code, "This route is not available")
	// })

	gg := app.Group("/v1/grants")
	gg.Get("/:value/lookup/:type", GrantsLookupRoute)
	gg.Get("/:value/lookup", GrantsLookupRoute)
	gg.Post("/:name/create", GrantsCreateRoute)
	// gg.RouteNotFound("/*", func(_ echo.Context) error {
	// 	return common.HTTPError(echo.ErrLocked.Code, "This route is not available")
	// })
}
