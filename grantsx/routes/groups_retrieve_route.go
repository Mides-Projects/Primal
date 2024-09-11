package routes

import (
	"github.com/holypvp/primal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GroupsRetrieveRoute(c echo.Context) error {
	values := service.Groups().All()

	body := make(map[string]map[string]interface{}, len(values))
	for _, g := range values {
		body[g.Id()] = g.Marshal()
	}

	return c.JSON(http.StatusOK, body)
}
