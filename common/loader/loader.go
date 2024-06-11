package loader

import (
	"github.com/gorilla/mux"
	server_routes "github.com/holypvp/primal/server/routes"
)

func LoadAll(router *mux.Router) {
	router.HandleFunc("/apiv2/servers/{id}/create/{port}", server_routes.ServerCreateRoute).Methods("POST")
	router.HandleFunc("/apiv2/servers/{id}/lookup", server_routes.LookupServers).Methods("GET")
	router.HandleFunc("/apiv2/servers/{id}/down", server_routes.ServerDownRoute).Methods("PATCH")
	router.HandleFunc("/apiv2/servers/{id}/tick", server_routes.ServerTickRoute).Methods("PATCH")

	router.HandleFunc("/apiv2/servers/groups/lookup", server_routes.GroupLookupRoute).Methods("GET")
}
