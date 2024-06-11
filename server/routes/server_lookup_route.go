package routes

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/holypvp/primal/common/middleware"
	"github.com/holypvp/primal/server"
	"log"
	"net/http"
)

func LookupServers(w http.ResponseWriter, r *http.Request) {
	if !middleware.HandleAuth(w, r) {
		return
	}

	serverId, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "No ID found", http.StatusBadRequest)

		return
	}

	err := json.NewEncoder(w).Encode(server.Service().Servers())
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)

		return
	}

	log.Print("Server " + serverId + " was looked up!")
}
