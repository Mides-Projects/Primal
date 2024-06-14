package routes

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/common/middleware"
	"github.com/holypvp/primal/server"
	"github.com/holypvp/primal/server/pubsub"
	"github.com/holypvp/primal/server/request"
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

	body := &request.ServerUpBody{}
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		log.Printf("[ServerUpRoute] Failed to parse request body: %v", err)

		return
	}

	serverInfo.SetDirectory(body.Directory)
	serverInfo.SetMotd(body.Motd)

	serverInfo.SetBungeeCord(body.BungeeCord)
	serverInfo.SetOnlineMode(body.OnlineMode)

	serverInfo.SetMaxSlots(body.MaxSlots)
	serverInfo.SetPlugins(body.Plugins)

	initialTime := serverInfo.InitialTime()
	if initialTime == 0 {
		http.Error(w, "Server was never down", http.StatusBadRequest)
		log.Printf("[ServerUpRoute] Server was never down")

		return
	}

	now := time.Now().UnixMilli()
	serverInfo.SetInitialTime(now)

	// Save the server model in a goroutine to avoid blocking the main thread
	go server.SaveModel(serverInfo.ToModel())

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

	log.Printf("[ServerUpRoute] Server %s is now back up. After %d ms", serverId, now-serverInfo.Heartbeat())
}
