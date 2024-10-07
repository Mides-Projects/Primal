package grants

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/model/grantsx"
	"github.com/holypvp/primal/model/player"
	"github.com/holypvp/primal/service"
	"net/http"
)

func LookupRoute(c fiber.Ctx) error {
	src := c.Params("src")
	if src == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'src' query parameter",
		})
	} else if src != "name" && src != "id" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid src",
		})
	} else if query := c.Query("query"); query == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'query' query parameter",
		})
	} else if !player.ValidQueryState(query) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "You can only lookup online or offline players",
		})
	} else if v := c.Params("value"); v == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'value' parameter",
		})
	} else if filter := c.Params("filter"); filter == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'filter' parameter",
		})
	} else if filter != "active" && filter != "all" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Parameter 'filter' must be either 'active', 'all' or empty",
		})
	} else {
		var (
			pi  *player.PlayerInfo
			err error
		)
		if src == "name" {
			pi, err = service.Player().UnsafeLookupByName(v, query == "online")
		} else {
			pi, err = service.Player().UnsafeLookupById(v, query == "online")
		}

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to lookup account: " + err.Error(),
			})
		} else if pi == nil && query == "online" {
			return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
				"message": "State is 'online' but account never joined",
			})
		} else if pi == nil {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "PlayerInfo not found",
			})
		}

		if query == "online" {
			service.Grants().Invalidate(pi.Id(), false)
		}

		tracker, err := service.Grants().UnsafeLookup(pi.Id(), query == "online")
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		} else if tracker == nil {
			tracker = grantsx.EmptyTracker()
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"id":     pi.Id(),
			"grants": tracker.Marshal(filter),
			"account": func() map[string]interface{} {
				if pi.Online() {
					return nil
				}

				return pi.Marshal()
			},
		})
	}
}
