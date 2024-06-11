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

func ServerUpRoute(w http.ResponseWriter, r *http.Request) {
	if !middleware.HandleAuth(w, r) {
		return
	}

	serverId, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "No ID found", http.StatusBadRequest)

		return
	}

	serverInfo := server.Service().LookupById(serverId)
	if serverInfo == nil {
		log.Println("Server not found... Failed to start server " + serverId)
	}

	log.Print("Server " + serverId + " was started!")

	payload, err := common.WrapPayload("API_SERVER_UP", pubsub.NewServerStatusPacket(serverId))
	if err != nil {
		http.Error(w, "Failed to marshal packet", http.StatusInternalServerError)

		return
	}

	err = common.RedisClient.Publish(context.Background(), common.RedisChannel, payload).Err()
	if err != nil {
		http.Error(w, "Failed to publish packet", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
