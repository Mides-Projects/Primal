package routes

import (
	"encoding/json"
	"github.com/holypvp/primal/common/middleware"
	"github.com/holypvp/primal/server"
	"net/http"
)

func LookupServers(w http.ResponseWriter, r *http.Request) {
	if !middleware.HandleAuth(w, r) {
		return
	}

	err := json.NewEncoder(w).Encode(server.Module().Values())
	if err == nil {
		return
	}

	http.Error(w, "Failed to encode response", http.StatusInternalServerError)
}
