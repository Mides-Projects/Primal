package account

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/service"
	"net/http"
	"sync/atomic"
)

func QuitRoute(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Missing 'id' parameter",
			"code":    http.StatusNotFound,
		})
	}

	acc := service.Account().LookupById(id)
	if acc == nil || !acc.Online() {
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"message": "You are not logged in",
			"code":    http.StatusServiceUnavailable,
		})
	}

	var body map[string]interface{}
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse body: " + err.Error(),
			"code":    http.StatusBadRequest,
		})
	}

	state := &atomic.Bool{}
	state.Store(true)
	defer acc.SetOnline(state.Load())

	displayName, ok := body["display_name"].(string)
	if !ok || displayName == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'display_name' parameter",
			"code":    http.StatusBadRequest,
		})
	}

	highestGroup, ok := body["highest_group"].(string)
	if !ok || highestGroup == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'highest_group' parameter",
			"code":    http.StatusBadRequest,
		})
	}

	timestamp, ok := body["timestamp"].(float64)
	if !ok || timestamp == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'timestamp' parameter",
			"code":    http.StatusBadRequest,
		})
	}

	canBeUpdated := false
	if acc.DisplayName() != displayName {
		acc.SetDisplayName(displayName)
		canBeUpdated = true
	}

	if acc.HighestGroup() != highestGroup {
		acc.SetHighestGroup(highestGroup)
		canBeUpdated = true
	}

	if canBeUpdated {
		go func() {
			if err := service.Account().Update(acc); err != nil {
				common.Log.Fatalf("Failed to update account: %v", err)
			}
		}()
	}

	if int64(timestamp) > acc.LastJoin().UnixMilli() {
		common.Log.Print("Player was disconnected due to a timestamp mismatch")
	} else {
		state.Store(false) // The state change to false because the player not was disconnected
	}

	// TODO: Broadcast a redis message to all servers that the player has logged out
	// with his display name and the server he was logged in

	// TODO: Cache the id account into a temporarily cache that will be deleted after a certain amount of time
	// this helps a lot to prevent make requests to database if they log in on less than 5 minutes

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "You have successfully logged out",
		"code":    http.StatusOK,
	})
}
