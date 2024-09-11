package startup

import (
	"context"
	"errors"
	"github.com/holypvp/primal/common"
	common_middleware "github.com/holypvp/primal/common/middleware"
	grantsx "github.com/holypvp/primal/grantsx/routes"
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
    common.Log = e.Logger

    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(common_middleware.HandleBasicAuth)

    db := common.MongoClient.Database("api")
    if err := loadGrantsX(e, db); err != nil {
        common.Log.Panicf("Failed to load grantsx: %v", err)
    }

    loadServerRoutes(e.Group("/v2/servers"), db)

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

func loadServerRoutes(g *echo.Group, db *mongo.Database) {
    // g.RouteNotFound("/*", func(c echo.Context) error {
    // 	return echo.NewHTTPError(echo.ErrLocked.Code, "This route is not available")
    // })
}

func loadGrantsX(e *echo.Echo, db *mongo.Database) error {
    g := e.Group("/v2/groups")
    g.POST("/:name/create/", grantsx.GroupCreateRoute)
    g.GET("/", grantsx.GroupsRetrieveRoute)
    g.RouteNotFound("/*", func(_ echo.Context) error {
        return common.HTTPError(echo.ErrLocked.Code, "This route is not available")
    })

    gg := e.Group("/v2/grants")
    gg.GET("/:value/lookup/:type", grantsx.GrantsLookupRoute)
    gg.GET("/:value/lookup", grantsx.GrantsLookupRoute)
    gg.POST("/:name/create", grantsx.GrantsCreateRoute)
    gg.RouteNotFound("/*", func(_ echo.Context) error {
        return common.HTTPError(echo.ErrLocked.Code, "This route is not available")
    })

    return nil
}

func Shutdown() {
    // TODO: Implement shutdown logic
}
