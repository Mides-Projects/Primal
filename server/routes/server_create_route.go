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
	"strconv"
	"time"
)

func ServerCreateRoute(w http.ResponseWriter, r *http.Request) {
	if !middleware.HandleAuth(w, r) {
		return
	}

	vars := mux.Vars(r)

	serverId, ok := vars["id"]
	if !ok {
		http.Error(w, "No ID found", http.StatusBadRequest)
		log.Print("[Server-Create] No ID found")

		return
	}

	serverInfo := server.Service().LookupById(serverId)
	if serverInfo != nil {
		http.Error(w, "Server already exists", http.StatusBadRequest)
		log.Printf("[Server-Create] Server %s already exists", serverId)

		return
	}

	port, ok := vars["port"]
	if !ok {
		http.Error(w, "No port found", http.StatusBadRequest)
		log.Print("[Server-Create] No port found")

		return
	}

	portNum, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		http.Error(w, "Invalid port", http.StatusBadRequest)
		log.Printf("[Server-Create] Invalid port: %v", err)

		return
	}

	if server.Service().LookupByPort(portNum) != nil {
		http.Error(w, "Port already in use", http.StatusBadRequest)
		log.Printf("[Server-Create] Port %d already in use", portNum)

		return
	}

	payload, err := common.WrapPayload("API_SERVER_CREATE", pubsub.NewServerCreatePacket(serverId, portNum))
	if err != nil {
		http.Error(w, "Failed to marshal packet", http.StatusInternalServerError)
		log.Printf("[Server-Create] Failed to marshal packet: %v", err)

		return
	}

	serverInfo = server.NewServerInfo(serverId, portNum)
	serverInfo.SetInitialTime(time.Now().UnixMilli())

	server.Service().AppendServer(serverInfo)

	// Save the model into MongoDB but in a goroutine, so it doesn't block the main thread
	go server.SaveModel(serverInfo.ToModel())

	err = common.RedisClient.Publish(context.Background(), common.RedisChannel, payload).Err()
	if err != nil {
		http.Error(w, "Failed to publish packet", http.StatusInternalServerError)
		log.Printf("[Server-Create] Failed to publish packet: %v", err)

		return
	}

	w.WriteHeader(http.StatusOK)

	log.Printf("[Server-Create] Server %s created on port %d", serverId, portNum)
}
