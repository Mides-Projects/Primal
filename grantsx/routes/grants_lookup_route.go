package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/service"
	"github.com/holypvp/primal/source/model"
	"net/http"
)

func GrantsLookupRoute(c fiber.Ctx) error {
	t := c.Params("type")
	if t != "" && t != "active" && t != "expired" {
		return common.HTTPError(c, http.StatusBadRequest, "Invalid type")
	}

	src := c.Query("src")
	if src == "" {
		return common.HTTPError(c, http.StatusBadRequest, "Missing 'src' query parameter")
	}

	if src != "name" && src != "id" {
		return common.HTTPError(c, http.StatusBadRequest, "No valid source found")
	}

	state := c.Query("state")
	if state == "" {
		return common.HTTPError(c, http.StatusBadRequest, "Missing 'state' query parameter")
	}

	if state != "online" && state != "offline" {
		return common.HTTPError(c, http.StatusBadRequest, "Invalid state")
	}

	v := c.Params("value")
	if v == "" {
		return common.HTTPError(c, http.StatusBadRequest, "Missing 'value' parameter")
	}

	var acc *model.Account
	var err error

	if src == "name" {
		acc, err = service.Account().UnsafeLookupByName(v)
	} else {
		acc, err = service.Account().UnsafeLookupById(v)
	}

	if err != nil {
		return common.HTTPError(c, http.StatusInternalServerError, "Failed to lookup account: "+err.Error())
	} else if acc == nil {
		return common.HTTPError(c, http.StatusNotFound, "Player not found")
	}

	ga, err := service.Grants().UnsafeLookup(acc.Id())
	if err != nil {
		return common.HTTPError(c, http.StatusInternalServerError, "Failed to lookup grant: "+err.Error())
	} else if ga == nil {
		return common.HTTPError(c, http.StatusNotFound, "Player not found")
	}

	if state == "online" && service.Grants().Lookup(ga.Account().Id()) == nil {
		service.Grants().Cache(ga)
	}

	return c.Status(http.StatusOK).JSON(marshalByType(t, ga.Marshal()))
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
