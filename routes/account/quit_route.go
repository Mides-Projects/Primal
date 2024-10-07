package account

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/model/player"
	"github.com/holypvp/primal/service"
	"net/http"
	"sync/atomic"
)

func QuitRoute(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Missing 'id' parameter",
		})
	}

	acc := service.Player().LookupById(id)
	if acc == nil || !acc.Online() {
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"message": "You are not logged in",
		})
	}

	var body player.UpdateBodyRequest
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse body: " + err.Error(),
		})
	}

	state := &atomic.Bool{}
	state.Store(true)
	defer acc.SetOnline(state.Load())

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

	canBeUpdated := false
	if acc.DisplayName() != body.DisplayName {
		acc.SetDisplayName(body.DisplayName)
		canBeUpdated = true
	}

	if acc.HighestGroup() != body.HighestGroup {
		acc.SetHighestGroup(body.HighestGroup)
		canBeUpdated = true
	}

	if canBeUpdated {
		go func() {
			if err := service.Player().Update(acc); err != nil {
				common.Log.Fatalf("Failed to update account: %v", err)
			}
		}()
	}

	if body.Timestamp > acc.LastJoin().UnixMilli() {
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
	})
}
