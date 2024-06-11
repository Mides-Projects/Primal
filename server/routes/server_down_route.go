package routes

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/common/middleware"
	"github.com/holypvp/primal/server"
	"github.com/holypvp/primal/server/pubsub"
	"net/http"
)

func ServerDownRoute(w http.ResponseWriter, r *http.Request) {
	if !middleware.HandleAuth(w, r) {
		return
	}

	id, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "No ID found", http.StatusBadRequest)

		return
	}

	serverInfo := server.Service().LookupById(id)
	if serverInfo == nil {
		http.Error(w, "Server not found", http.StatusNotFound)

		return
	}

	// TODO: Mark the server as down
	// TODO: Send a redis packet to all servers to update their server list

	result, err := json.Marshal(pubsub.NewServerDownPacket(serverInfo.Id()))
	if err != nil {
		http.Error(w, "Failed to marshal packet", http.StatusInternalServerError)

		return
	}

	err = common.RedisClient.Publish(context.Background(), "apiv2", result).Err()
	if err != nil {
		http.Error(w, "Failed to publish packet", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
