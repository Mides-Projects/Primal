package startup

import (
	"context"
	"errors"
	common_middleware "github.com/holypvp/primal/common/middleware"
	server_routes "github.com/holypvp/primal/server/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func LoadAll(now time.Time, port string) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(common_middleware.HandleBasicAuth)

	loadServerRoutes(e.Group("/v2/servers"))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		e.Logger.Printf("App take %s to start\n", time.Since(now))
		if err := e.Start(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
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
