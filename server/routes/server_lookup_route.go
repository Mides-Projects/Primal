package routes

import (
	"github.com/holypvp/primal/server"
	"github.com/holypvp/primal/server/response"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func LookupServers(c echo.Context) error {
	serverId := c.Param("id")
	if serverId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "No ID found")
	}

	var responses []response.ServerInfoResponse
	for _, serverInfo := range server.Service().Servers() {
		responses = append(responses, response.NewServerInfoResponse(serverInfo))
	}

	if responses == nil {
		responses = []response.ServerInfoResponse{}
	}

	log.Print("Server " + serverId + " has requested all servers")

	return c.JSON(http.StatusOK, responses)
}
