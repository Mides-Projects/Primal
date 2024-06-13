package group

import (
	"encoding/json"
	"github.com/holypvp/primal/common/middleware"
	"github.com/holypvp/primal/server"
	"github.com/holypvp/primal/server/response"
	"log"
	"net/http"
)

func GroupLookupRoute(w http.ResponseWriter, r *http.Request) {
	if !middleware.HandleAuth(w, r) {
		return
	}

	// serverId, ok := mux.Vars(r)["id"]
	// if !ok {
	// 	http.Error(w, "No ID found", http.StatusBadRequest)
	//
	// 	return
	// }
	//
	// serverInfo := server.Service().LookupById(serverId)
	// if serverInfo == nil {
	// 	http.Error(w, "Server not found", http.StatusNotFound)
	//
	// 	return
	// }

	var groupsResponse []response.ServerGroupResponse
	for _, group := range server.Service().Groups() {
		groupsResponse = append(groupsResponse, response.ServerGroupResponse{
			Id:                    group.Id(),
			Metadata:              group.Metadata(),
			Announcements:         group.Announcements(),
			AnnouncementsInterval: group.AnnouncementsInterval(),
			FallbackServerId:      group.FallbackServerId(),
		})
	}

	if groupsResponse == nil {
		groupsResponse = []response.ServerGroupResponse{}
	}

	err := json.NewEncoder(w).Encode(groupsResponse)
	if err != nil {
		http.Error(w, "Failed to marshal groups", http.StatusInternalServerError)

		return
	}

	log.Print("Server group lookup route hit")
}
