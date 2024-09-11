package routes

import (
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GrantsLookupRoute(e echo.Context) error {
	t := e.Param("type")
	if t != "" && t != "active" && t != "expired" {
		return common.HTTPError(http.StatusBadRequest, "Invalid type")
	}

	src := e.QueryParam("src")
	if src == "" {
		return common.HTTPError(http.StatusBadRequest, "Missing 'src' query parameter")
	}

	if src != "name" && src != "id" {
		return common.HTTPError(http.StatusBadRequest, "No valid source found")
	}

	state := e.QueryParam("state")
	if state == "" {
		return common.HTTPError(http.StatusBadRequest, "Missing 'state' query parameter")
	}

	if state != "online" && state != "offline" {
		return common.HTTPError(http.StatusBadRequest, "Invalid state")
	}

	v := e.Param("value")
	if src == "name" {
		v = service.Account().FetchAccountId(v)
	}

	if v == "" {
		return common.HTTPError(http.StatusBadRequest, "Missing 'id' into our database")
	}

	ga, err := service.Grants().Lookup(v)
	if err != nil {
		return common.HTTPError(http.StatusInternalServerError, "Failed to lookup grant: "+err.Error())
	} else if ga == nil {
		return common.HTTPError(http.StatusNotFound, "Player not found")
	}

	if state == "online" && service.Grants().LookupAtCache(ga.Account().Id()) == nil {
		service.Grants().Cache(ga)
	}

	return e.JSON(http.StatusOK, marshalByType(t, ga.Marshal()))
}

func marshalByType(t string, body map[string]interface{}) map[string]interface{} {
	if t == "" {
		return body
	}

	if t == "active" {
		delete(body, "expired_grants")
	} else {
		delete(body, "active_grants")
	}

	return body
}
