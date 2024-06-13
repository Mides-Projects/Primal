package routes

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/holypvp/primal/common/middleware"
	"github.com/holypvp/primal/server"
	"github.com/holypvp/primal/server/response"
	"log"
	"net/http"
)

func LookupServers(w http.ResponseWriter, r *http.Request) {
	if !middleware.HandleAuth(w, r) {
		return
	}

	serverId, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "No ID found", http.StatusBadRequest)
		log.Printf("[Server-Lookup] No ID found")

		return
	}

	var responses []response.ServerInfoResponse
	for _, serverInfo := range server.Service().Servers() {
		responses = append(responses, response.ServerInfoResponse{
			Id:             serverInfo.Id(),
			Port:           serverInfo.Port(),
			Groups:         serverInfo.Groups(),
			PlayersCount:   serverInfo.PlayersCount(),
			MaxSlots:       serverInfo.MaxSlots(),
			Heartbeat:      serverInfo.Heartbeat(),
			Players:        serverInfo.Players(),
			BungeeCord:     serverInfo.BungeeCord(),
			OnlineMode:     serverInfo.OnlineMode(),
			ActiveThreads:  serverInfo.ActiveThreads(),
			DaemonThreads:  serverInfo.DaemonThreads(),
			Motd:           serverInfo.Motd(),
			TicksPerSecond: serverInfo.TicksPerSecond(),
			Directory:      serverInfo.Directory(),
			FullTicks:      serverInfo.FullTicks(),
			InitialTime:    serverInfo.InitialTime(),
			Plugins:        serverInfo.Plugins(),
		})
	}

	if responses == nil {
		responses = []response.ServerInfoResponse{}
	}

	err := json.NewEncoder(w).Encode(responses)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("[Server-Lookup] Failed to encode response: %v", err)

		return
	}

	log.Print("Server " + serverId + " was looked up!")
}
