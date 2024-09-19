package account

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/service"
	"net/http"
)

func UpdateRoute(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'id' parameter",
			"code":    http.StatusBadRequest,
		})
	}

	var body map[string]interface{}
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to bind request body: " + err.Error(),
			"code":    http.StatusBadRequest,
		})
	}

	displayName, ok := body["display_name"].(string)
	if !ok || displayName == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'display_name' field",
			"code":    http.StatusBadRequest,
		})
	}

	operator, ok := body["operator"].(bool)
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'operator' field",
			"code":    http.StatusBadRequest,
		})
	}

	highestGroup, ok := body["highest_group"].(string)
	if !ok || highestGroup == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'highest_group' field",
			"code":    http.StatusBadRequest,
		})
	}

	go func() {
		if acc, err := service.Account().UnsafeLookupById(id); err != nil {
			common.Log.Fatalf("Failed to lookup account: %s", err)
		} else if acc == nil {
			common.Log.Fatalf("Account not found: %s", id)
		} else {
			canBeSaved := false
			if acc.DisplayName() != displayName {
				acc.SetDisplayName(displayName)
				canBeSaved = true
			}

			if acc.Operator() != operator {
				acc.SetOperator(operator)
				canBeSaved = true
			}

			if acc.HighestGroup() != highestGroup {
				acc.SetHighestGroup(highestGroup)
				canBeSaved = true
			}

			if !canBeSaved {
				return
			}

			if err = service.Account().Update(acc); err != nil {
				common.Log.Fatalf("Failed to update account: %s", err)
			}
		}
	}()

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Account has been updated",
		"code":    http.StatusOK,
	})
}
