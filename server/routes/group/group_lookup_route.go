package group

import (
	"github.com/holypvp/primal/server/response"
	"github.com/holypvp/primal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GroupLookupRoute(c echo.Context) error {
	// serverId, ok := mux.Vars(r)["id"]
	// if !ok {
	// 	http.Error(w, "No ID found", http.StatusBadRequest)
	//
	// 	return
	// }
	//
	// serverInfo := server.Server().LookupById(serverId)
	// if serverInfo == nil {
	// 	http.Error(w, "Server not found", http.StatusNotFound)
	//
	// 	return
	// }

	var groupsResponse []response.ServerGroupResponse
	for _, group := range service.Server().Groups() {
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

	return c.JSON(http.StatusOK, groupsResponse)
}
