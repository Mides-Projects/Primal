package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/server/response"
	"github.com/holypvp/primal/service"
	"log"
	"net/http"
)

func LookupServers(c fiber.Ctx) error {
	serverId := c.Params("id")
	if serverId == "" {
		return common.HTTPError(c, http.StatusBadRequest, "No server ID found")
	}

	var responses []response.ServerInfoResponse
	for _, serverInfo := range service.Server().Servers() {
		responses = append(responses, response.NewServerInfoResponse(serverInfo))
	}

	if responses == nil {
		responses = []response.ServerInfoResponse{}
	}

	log.Print("Server " + serverId + " has requested all servers")

	return c.Status(http.StatusOK).JSON(responses)
}
