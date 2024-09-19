package grants

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/account"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/service"
	"net/http"
)

func GrantsLookupRoute(c fiber.Ctx) error {
	t := c.Params("type")
	if t != "" && t != "active" && t != "expired" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid type",
			"code":    http.StatusBadRequest,
		})
	}

	src := c.Query("src")
	if src == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'src' query parameter",
			"code":    http.StatusBadRequest,
		})
	}

	if src != "name" && src != "id" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid src",
			"code":    http.StatusBadRequest,
		})
	}

	state := c.Query("state")
	if state == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'state' query parameter",
			"code":    http.StatusBadRequest,
		})
	}

	if state != "online" && state != "offline" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid state",
			"code":    http.StatusBadRequest,
		})
	}

	v := c.Params("value")
	if v == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'value' parameter",
			"code":    http.StatusBadRequest,
		})
	}

	var acc *account.Account
	var err error

	if state == "online" {
		if src == "name" {
			acc = service.Account().LookupByName(v)
		} else {
			acc = service.Account().LookupById(v)
		}
	} else {
		if src == "name" {
			acc, err = service.Account().UnsafeLookupByName(v)
		} else {
			acc, err = service.Account().UnsafeLookupById(v)
		}
	}

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to lookup account: " + err.Error(),
			"code":    http.StatusInternalServerError,
		})
	} else if acc == nil && state == "online" {
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"message": "Account is not online, but the state is set to online",
			"code":    http.StatusServiceUnavailable,
		})
	} else if acc == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Account not found",
			"code":    http.StatusNotFound,
		})
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
