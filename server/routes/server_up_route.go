package routes

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/common/middleware"
	"github.com/holypvp/primal/server"
	"github.com/holypvp/primal/server/model"
	"github.com/holypvp/primal/server/pubsub"
	"log"
	"net/http"
	"time"
)

func ServerUpRoute(w http.ResponseWriter, r *http.Request) {
	if !middleware.HandleAuth(w, r) {
		return
	}

	serverId, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "No ID found", http.StatusBadRequest)
		log.Printf("[ServerUpRoute] No ID found")

		return
	}

	serverInfo := server.Service().LookupById(serverId)
	if serverInfo == nil {
		http.Error(w, "Server not found", http.StatusBadRequest)
		log.Printf("[ServerUpRoute] Server not found")

		return
	}

	bungeeMode := r.URL.Query().Get("bungee")
	if bungeeMode == "" {
		http.Error(w, "Bungee mode is required", http.StatusBadRequest)
		log.Printf("[ServerUpRoute] Bungee mode is required")

		return
	}

	if bungeeMode != "true" && bungeeMode != "false" {
		http.Error(w, "Invalid bungee mode", http.StatusBadRequest)
		log.Printf("[ServerUpRoute] Invalid bungee mode")

		return
	}

	onlineMode := r.URL.Query().Get("online")
	if onlineMode == "" {
		http.Error(w, "Online mode is required", http.StatusBadRequest)
		log.Printf("[ServerUpRoute] Online mode is required")

		return
	}

	if onlineMode != "true" && onlineMode != "false" {
		http.Error(w, "Invalid online mode", http.StatusBadRequest)
		log.Printf("[ServerUpRoute] Invalid online mode")

		return
	}

	initialTime := serverInfo.InitialTime()
	if initialTime == 0 {
		http.Error(w, "Server was never down", http.StatusBadRequest)
		log.Printf("[ServerUpRoute] Server was never down")

		return
	}

	serverInfo.SetBungeeCord(bungeeMode == "true")
	serverInfo.SetOnlineMode(onlineMode == "true")

	now := time.Now().UnixMilli()
	serverInfo.SetInitialTime(now)

	// Save the server model in a goroutine to avoid blocking the main thread
	go server.SaveModel(model.WrapServerInfo(serverInfo))

	payload, err := common.WrapPayload("API_SERVER_UP", pubsub.NewServerStatusPacket(serverId))
	if err != nil {
		http.Error(w, "Failed to marshal packet", http.StatusInternalServerError)
		log.Printf("[ServerUpRoute] Failed to marshal packet: %v", err)

		return
	}

	err = common.RedisClient.Publish(context.Background(), common.RedisChannel, payload).Err()
	if err != nil {
		http.Error(w, "Failed to publish packet", http.StatusInternalServerError)
		log.Printf("[ServerUpRoute] Failed to publish packet: %v", err)

		return
	}

	w.WriteHeader(http.StatusOK)

	log.Printf("[ServerUpRoute] Server %s is now back up. After %d ms", serverId, now-initialTime)
}
