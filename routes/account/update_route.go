package account

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/model/player"
	"github.com/holypvp/primal/service"
	"net/http"
)

func UpdateRoute(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'id' parameter",
		})
	}

	var body player.UpdateBodyRequest
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to bind request body: " + err.Error(),
		})
	}

	if body.DisplayName == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'display_name' body field",
		})
	}

	if body.HighestGroup == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'highest_group' body field",
		})
	}

	if body.Timestamp == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'timestamp' body field",
		})
	}

	go func() {
		if pi, err := service.Player().UnsafeLookupById(id, false); err != nil {
			common.Log.Fatalf("Failed to lookup account: %s", err)
		} else if pi == nil {
			common.Log.Fatalf("Player not found: %s", id)
		} else {
			canBeSaved := false
			if pi.DisplayName() != body.DisplayName {
				pi.SetDisplayName(body.DisplayName)
				canBeSaved = true
			}

			if pi.Operator() != body.Operator {
				pi.SetOperator(body.Operator)
				canBeSaved = true
			}

			if pi.HighestGroup() != body.HighestGroup {
				pi.SetHighestGroup(body.HighestGroup)
				canBeSaved = true
			}

			if !canBeSaved {
				return
			}

			if err = service.Player().Update(pi); err != nil {
				common.Log.Fatalf("Failed to update account: %s", err)
			}
		}
	}()

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Player has been updated",
	})
}
