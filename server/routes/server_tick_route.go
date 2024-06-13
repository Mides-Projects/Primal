package routes

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/common/middleware"
	"github.com/holypvp/primal/server"
	"github.com/holypvp/primal/server/request"
	"log"
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

	serverInfo := server.Service().LookupById(id)
	if serverInfo == nil {
		http.Error(w, "Server not found", http.StatusNotFound)

		return
	}

	body := &request.ServerTickBody{}
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)

		return
	}

	// serverInfo.SetGroups(body.Groups)
	serverInfo.SetPlayersCount(body.PlayersCount)
	// serverInfo.SetPort(body.Port)

	serverInfo.SetHeartbeat(body.Heartbeat)
	// serverInfo.SetBungeeCord(body.BungeeCord)
	// serverInfo.SetOnlineMode(body.OnlineMode)

	serverInfo.SetActiveThreads(body.ActiveThreads)
	serverInfo.SetDaemonThreads(body.DaemonThreads)

	// serverInfo.SetMotd(body.Motd)
	serverInfo.SetTicksPerSecond(body.TicksPerSecond)
	// serverInfo.SetDirectory(body.Directory)
	serverInfo.SetFullTicks(body.FullTicks)
	// serverInfo.SetMaxSlots(body.MaxSlots)
	// serverInfo.SetInitialTime(body.InitialTime)
	// serverInfo.SetPlugins(body.Plugins)
	serverInfo.SetPlayers(body.Players)

	payload, err := common.WrapPayload("API_SERVER_TICK", body)
	if err != nil {
		http.Error(w, "Failed to marshal packet", http.StatusInternalServerError)

		return
	}

	err = common.RedisClient.Publish(context.Background(), common.RedisChannel, payload).Err()
	if err != nil {
		http.Error(w, "Failed to publish packet", http.StatusInternalServerError)

		return
	}

	log.Printf("Successfully updated server tick for %s\n", id)

	w.WriteHeader(http.StatusOK)
}
