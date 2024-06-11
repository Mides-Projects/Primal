package routes

import (
	"encoding/json"
	"github.com/holypvp/primal/common/middleware"
	"github.com/holypvp/primal/server"
	"net/http"
)

func GroupLookupRoute(w http.ResponseWriter, r *http.Request) {
	if !middleware.HandleAuth(w, r) {
		return
	}

	// serverId, ok := mux.Vars(r)["id"]
	// if !ok {
	// 	http.Error(w, "No ID found", http.StatusBadRequest)
	//
	// 	return
	// }
	//
	// serverInfo := server.Service().LookupById(serverId)
	// if serverInfo == nil {
	// 	http.Error(w, "Server not found", http.StatusNotFound)
	//
	// 	return
	// }

	err := json.NewEncoder(w).Encode(server.Service().Groups())
	if err != nil {
		http.Error(w, "Failed to marshal groups", http.StatusInternalServerError)

		return
	}
}
