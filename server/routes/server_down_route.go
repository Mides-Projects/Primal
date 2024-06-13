package routes

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/common/middleware"
	"github.com/holypvp/primal/server"
	"github.com/holypvp/primal/server/pubsub"
	"log"
	"net/http"
)

func ServerDownRoute(w http.ResponseWriter, r *http.Request) {
	if !middleware.HandleAuth(w, r) {
		return
	}

	id, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "No ID found", http.StatusBadRequest)
		log.Print("[Server-Down] No ID found")

		return
	}

	serverInfo := server.Service().LookupById(id)
	if serverInfo == nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		log.Print("[Server-Down] Server not found")

		return
	}

	payload, err := common.WrapPayload("API_SERVER_DOWN", pubsub.NewServerStatusPacket(serverInfo.Id()))
	if err != nil {
		http.Error(w, "Failed to marshal packet", http.StatusInternalServerError)
		log.Printf("[Server-Down] Failed to marshal packet: %v", err)

		return
	}

	err = common.RedisClient.Publish(context.Background(), common.RedisChannel, payload).Err()
	if err != nil {
		http.Error(w, "Failed to publish packet", http.StatusInternalServerError)
		log.Printf("[Server-Down] Failed to publish packet: %v", err)

		return
	}

	w.WriteHeader(http.StatusOK)

	log.Printf("[Server-Down] Server %s is now down", serverInfo.Id())
}
