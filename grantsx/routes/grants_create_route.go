package routes

import (
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GrantsCreateRoute(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return common.HTTPError(http.StatusBadRequest, "No name found for the player")
	}

	id := service.Account().FetchAccountId(name)
	if id == "" {
		return common.HTTPError(http.StatusNotFound, "Player not found")
	}

	var body map[string]interface{}
	if err := c.Bind(&body); err != nil {
		return common.HTTPError(http.StatusBadRequest, "Failed to bind request body: "+err.Error())
	}

	// TODO: Here apply the grant to the player

	return nil
}
