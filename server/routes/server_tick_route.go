package routes

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/holypvp/primal/common/middleware"
	"github.com/holypvp/primal/server"
	"github.com/holypvp/primal/server/response"
	"net/http"
)

func ServerTickRoute(w http.ResponseWriter, r *http.Request) {
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

	result := &response.ServerTickRequest{}
	err := json.NewDecoder(r.Body).Decode(result)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)

		return
	}

	serverInfo.ClearGroups()

	for _, group := range result.Groups {
		serverInfo.AddGroup(group)
	}

	w.WriteHeader(http.StatusOK)
}
