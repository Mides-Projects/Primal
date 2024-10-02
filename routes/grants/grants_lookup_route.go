package grants

import (
    "github.com/gofiber/fiber/v3"
    "github.com/holypvp/primal/model"
    "github.com/holypvp/primal/model/grantsx"
    "github.com/holypvp/primal/service"
    "net/http"
)

func LookupRoute(c fiber.Ctx) error {
    filter := c.Params("filter")
    if filter == "" {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "message": "Missing 'filter' parameter",
        })
    }

    if filter != "active" && filter != "all" {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "message": "Parameter 'filter' must be either 'active', 'all' or empty",
        })
    }

    src := c.Query("src")
    if src == "" {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "message": "Missing 'src' query parameter",
        })
    }

    if src != "name" && src != "id" {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "message": "Invalid src",
        })
    }

    state := c.Query("state")
    if state == "" {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "message": "Missing 'state' query parameter",
        })
    }

    onlineState := state == "online"
    if !onlineState && state != "offline" {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "message": "You can only lookup online or offline players",
        })
    }

    v := c.Params("value")
    if v == "" {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "message": "Missing 'value' parameter",
        })
    }

    var (
        acc *model.Account
        err error
    )
    if src == "name" {
        acc, err = service.Account().UnsafeLookupByName(v, onlineState)
    } else {
        acc, err = service.Account().UnsafeLookupById(v, onlineState)
    }

    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to lookup account: " + err.Error(),
        })
    } else if acc == nil && onlineState {
        return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
            "message": "State is 'online' but account never joined",
        })
    } else if acc == nil {
        return c.Status(http.StatusNotFound).JSON(fiber.Map{
            "message": "Account not found",
        })
    }

    if onlineState {
        service.Grants().Invalidate(acc.Id())
    }

    tracker, err := service.Grants().UnsafeLookup(acc.Id(), state == "online")
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "message": err.Error(),
        })
    } else if tracker == nil {
        tracker = grantsx.EmptyTracker()
    }

    return c.Status(http.StatusOK).JSON(fiber.Map{
        "id":     acc.Id(),
        "grants": tracker.Marshal(filter),
        "account": func() map[string]interface{} {
            if acc.Online() {
                return nil
            }

            return acc.Marshal()
        },
    })
}
