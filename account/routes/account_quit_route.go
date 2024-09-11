package routes

import (
    "github.com/gofiber/fiber/v3"
    "github.com/holypvp/primal/service"
    "net/http"
)

func AccountQuitRoute(c fiber.Ctx) error {
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

    acc.SetOnline(false)

    // TODO: Broadcast a redis message to all servers that the player has logged out
    // with his display name and the server he was logged in

    // TODO: Cache the id account into a temporarily cache that will be deleted after a certain amount of time
    // this helps a lot to prevent make requests to database if they log in on less than 5 minutes

    return c.Status(http.StatusOK).JSON(fiber.Map{
        "message": "You have successfully logged out",
        "code":    http.StatusOK,
    })
}
