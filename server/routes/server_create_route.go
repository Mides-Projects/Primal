package routes

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/common/middleware"
	"github.com/holypvp/primal/server"
	"github.com/holypvp/primal/server/pubsub"
	"net/http"
	"strconv"
)

func ServerCreateRoute(w http.ResponseWriter, r *http.Request) {
	if !middleware.HandleAuth(w, r) {
		return
	}

	vars := mux.Vars(r)

	serverId, ok := vars["id"]
	if !ok {
		http.Error(w, "No ID found", http.StatusBadRequest)

		return
	}

	serverInfo := server.Service().LookupById(serverId)
	if serverInfo != nil {
		http.Error(w, "Server already exists", http.StatusBadRequest)

		return
	}

	port, ok := vars["port"]
	if !ok {
		http.Error(w, "No port found", http.StatusBadRequest)

		return
	}

	portNum, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		http.Error(w, "Invalid port", http.StatusBadRequest)

		return
	}

	if server.Service().LookupByPort(portNum) != nil {
		http.Error(w, "Port already in use", http.StatusBadRequest)

		return
	}

	payload, err := common.WrapPayload("API_SERVER_CREATE", pubsub.NewServerCreatePacket(serverId, portNum))
	if err != nil {
		http.Error(w, "Failed to marshal packet", http.StatusInternalServerError)

		return
	}

	err = common.RedisClient.Publish(context.Background(), "apiv2", payload).Err()
	if err != nil {
		http.Error(w, "Failed to publish packet", http.StatusInternalServerError)

		return
	}

	serverInfo = server.NewServerInfo(serverId, portNum)
	server.Service().AppendServer(serverInfo)

	w.WriteHeader(http.StatusOK)
}
