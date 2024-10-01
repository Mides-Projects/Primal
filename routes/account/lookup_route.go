package account

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/model"
	"github.com/holypvp/primal/service"
	"net/http"
)

func LookupRoute(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'id' parameter",
		})
	} else if src := c.Query("src"); src == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'src' query parameter",
		})
	} else if src != "name" && src != "id" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid 'src' query parameter",
		})
	} else {
		var (
			acc *model.Account
			err error
		)
		if src == "name" {
			acc, err = service.Account().UnsafeLookupByName(id)
		} else {
			acc, err = service.Account().UnsafeLookupById(id)
		}

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to lookup account: " + err.Error(),
			})
		}

		if acc == nil {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "Player not found",
			})
		}

		return c.Status(http.StatusOK).JSON(acc.Marshal())
	}
}
