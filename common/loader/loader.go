package loader

import (
	common_middleware "github.com/holypvp/primal/common/middleware"
	server_routes "github.com/holypvp/primal/server/routes"
	"github.com/holypvp/primal/server/routes/group"
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

	e.POST("/apiv2/servers/:id/create/:port", server_routes.ServerCreateRoute)
	e.GET("/apiv2/servers/:id/lookup", server_routes.LookupServers)
	e.PATCH("/apiv2/servers/:id/down", server_routes.ServerDownRoute)
	e.PATCH("/apiv2/servers/:id/up", server_routes.ServerUpRoute)
	e.PATCH("/apiv2/servers/:id/tick", server_routes.ServerTickRoute)

	e.GET("/apiv2/servers/:id/groups/lookup", group.GroupLookupRoute)

	log.Printf("App take %s to start\n", time.Since(now))
	log.Fatal(e.Start("0.0.0.0:" + port))

	// router.HandleFunc("/apiv2/servers/{id}/create/{port}", server_routes.ServerCreateRoute).Methods("POST")
	// router.HandleFunc("/apiv2/servers/{id}/lookup", server_routes.LookupServers).Methods("GET")
	// router.HandleFunc("/apiv2/servers/{id}/down", server_routes.ServerDownRoute).Methods("PATCH")
	// router.HandleFunc("/apiv2/servers/{id}/up", server_routes.ServerUpRoute).Methods("PATCH")
	// router.HandleFunc("/apiv2/servers/{id}/tick", server_routes.ServerTickRoute).Methods("PATCH")
	//
	// router.HandleFunc("/apiv2/servers/{id}/groups/lookup", group.GroupLookupRoute).Methods("GET")
}
