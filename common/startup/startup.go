package startup

import (
	common_middleware "github.com/holypvp/primal/common/middleware"
	server_routes "github.com/holypvp/primal/server/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"time"
)

func LoadAll(now time.Time, port string) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(common_middleware.HandleBasicAuth)

	loadServerRoutes(e.Group("/v2/servers"))

	log.Printf("App take %s to start\n", time.Since(now))
	log.Fatal(e.Start("0.0.0.0:" + port))
}

func loadServerRoutes(g *echo.Group) {
	g.POST("/:id/create/:port", server_routes.ServerCreateRoute)
	g.GET("/:id/lookup", server_routes.LookupServers)
	g.PATCH("/:id/down", server_routes.ServerDownRoute)
	g.PATCH("/:id/up", server_routes.ServerUpRoute)
	g.PATCH("/:id/tick", server_routes.ServerTickRoute)

	g.RouteNotFound("/*", func(c echo.Context) error {
		return echo.NewHTTPError(echo.ErrLocked.Code, "This route is not available")
	})
}

func Shutdown() {
	// TODO: Implement shutdown logic
}
