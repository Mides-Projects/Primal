package loader

import (
	"github.com/gorilla/mux"
	server_routes "github.com/holypvp/primal/server/routes"
)

func LoadAll(router *mux.Router) {
	router.HandleFunc("/api/v2/servers/lookup", server_routes.LookupServers).Methods("GET")
	router.HandleFunc("/api/v2/servers/{id}/down", server_routes.ServerDownRoute).Methods("POST")
	router.HandleFunc("/api/v2/servers/{id}/tick", server_routes.ServerTickRoute).Methods("PATCH")
}
