package group

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/server/response"
	"github.com/holypvp/primal/service"
	"net/http"
)

func GroupLookupRoute(c fiber.Ctx) error {
	// serverId, ok := mux.Vars(r)["id"]
	// if !ok {
	// 	http.Error(w, "No ID found", http.StatusBadRequest)
	//
	// 	return
	// }
	//
	// serverInfo := server.Server().UnsafeLookupById(serverId)
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

	return c.Status(http.StatusOK).JSON(groupsResponse)
}
