package startup

import (
	"context"
	"errors"
	"github.com/holypvp/primal/common"
	common_middleware "github.com/holypvp/primal/common/middleware"
	"github.com/holypvp/primal/server"
	server_routes "github.com/holypvp/primal/server/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
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

	loadServerRoutes(e.Group("/v2/servers"), common.MongoClient.Database("server-monitor"))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	common.Log = e.Logger

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

func loadServerRoutes(g *echo.Group, db *mongo.Database) {
	if err := server.LoadServers(db); err != nil {
		common.Log.Panicf("Failed to load servers: %v", err)
	}
	if err := server.LoadGroups(db); err != nil {
		common.Log.Panicf("Failed to load server groups: %v", err)
	}

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
