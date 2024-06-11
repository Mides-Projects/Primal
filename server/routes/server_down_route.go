package routes

import (
	"github.com/gorilla/mux"
	"github.com/holypvp/primal/common/middleware"
	"github.com/holypvp/primal/server"
	"net/http"
)

func ServerDownRoute(w http.ResponseWriter, r *http.Request) {
	if !middleware.HandleAuth(w, r) {
		return
	}

	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		http.Error(w, "No ID found", http.StatusBadRequest)

		return
	}

	serverInfo := server.Module().LookupById(id)
	if serverInfo == nil {
		http.Error(w, "Server not found", http.StatusNotFound)

		return
	}

	// TODO: Mark the server as down
	// TODO: Send a redis packet to all servers to update their server list

	w.WriteHeader(http.StatusOK)
}
