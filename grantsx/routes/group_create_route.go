package routes

import (
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/grantsx/model"
	"github.com/holypvp/primal/grantsx/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GroupCreateRoute(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return common.HTTPError(http.StatusBadRequest, "No group name found")
	}

	if service.Groups().LookupByName(name) != nil {
		return common.HTTPError(http.StatusConflict, "Group already exists")
	}

	g := model.EmptyGroup(name)
	service.Groups().Cache(g)

	go func() {
		if err := service.Groups().Save(g); err != nil {
			common.Log.Errorf("Failed to save group: %v", err)
		}
	}()

	return c.JSON(http.StatusOK, g)
}
