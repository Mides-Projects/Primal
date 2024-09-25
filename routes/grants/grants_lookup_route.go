package grants

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/model"
	"github.com/holypvp/primal/model/grantsx"
	"github.com/holypvp/primal/service"
	"net/http"
)

func LookupRoute(c fiber.Ctx) error {
	t := c.Params("type")
	if t != "" && t != "active" && t != "expired" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Parameter 'type' must be either 'active', 'expired' or empty",
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
			"message": "You can only lookup online or offline players",
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

	var (
		acc *model.Account
		err error
	)
	if src == "name" {
		acc, err = service.Account().UnsafeLookupByName(v)
	} else {
		acc, err = service.Account().UnsafeLookupById(v)
	}

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"code":    http.StatusInternalServerError,
		})
	} else if acc == nil && state == "online" {
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"message": "State is 'online' but account never joined",
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
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"code":    http.StatusInternalServerError,
		})
	} else if ga == nil {
		ga = grantsx.EmptyGrantsAccount(acc)
	}

	if state == "online" && service.Grants().Lookup(ga.Account().Id()) == nil {
		service.Grants().Cache(ga)
	}

	return c.Status(http.StatusOK).JSON(marshal(ga, t))
}

// marshal returns the grantsx account as a map.
func marshal(ga *grantsx.Tracker, t string) map[string]interface{} {
	body := map[string]interface{}{}
	if t == "" {
		body["expired"] = ga.ExpiredGrants()
	}

	body["active"] = ga.ActiveGrants()
	body["id"] = ga.Account().Id()

	var all []map[string]interface{}
	for _, g := range ga.ActiveGrants() {
		all = append(all, g.Marshal())
	}

	if t != "active" {
		for _, g := range ga.ExpiredGrants() {
			all = append(all, g.Marshal())
		}
	}

	body["all"] = all

	return body
}
