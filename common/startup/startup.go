package startup

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func LoadAll(now time.Time, port string) {

}

func loadServerRoutes(g *echo.Group, db *mongo.Database) {
	// g.RouteNotFound("/*", func(c echo.Context) error {
	// 	return echo.NewHTTPError(echo.ErrLocked.Code, "This route is not available")
	// })
}

func Shutdown() {
	// TODO: Implement shutdown logic
}
