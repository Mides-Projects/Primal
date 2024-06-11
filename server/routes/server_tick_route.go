package routes

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/holypvp/primal/common"
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

	serverInfo := server.Service().LookupById(id)
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

	serverInfo.SetGroups(result.Groups)
	serverInfo.SetPlayersCount(result.PlayersCount)
	serverInfo.SetPort(result.Port)

	serverInfo.SetHeartbeat(result.Heartbeat)
	serverInfo.SetBungeeCord(result.BungeeCord)
	serverInfo.SetOnlineMode(result.OnlineMode)

	serverInfo.SetActiveThreads(result.ActiveThreads)
	serverInfo.SetDaemonThreads(result.DaemonThreads)

	serverInfo.SetMotd(result.Motd)
	serverInfo.SetTicksPerSecond(result.TicksPerSecond)
	serverInfo.SetDirectory(result.Directory)
	serverInfo.SetFullTicks(result.FullTicks)
	serverInfo.SetMaxSlots(result.MaxSlots)
	serverInfo.SetInitialTime(result.InitialTime)
	serverInfo.SetPlugins(result.Plugins)
	serverInfo.SetPlayers(result.Players)

	payload, err := common.WrapPayload("API_SERVER_TICK", result)
	if err != nil {
		http.Error(w, "Failed to marshal packet", http.StatusInternalServerError)

		return
	}

	err = common.RedisClient.Publish(context.Background(), "apiv2", payload).Err()
	if err != nil {
		http.Error(w, "Failed to publish packet", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
