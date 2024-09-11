package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/account"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/service"
	"net/http"
)

func AccountJoinRoute(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return common.HTTPError(c, http.StatusBadRequest, "No account ID found")
	}

	name := c.Query("name")
	if name == "" {
		return common.HTTPError(c, http.StatusBadRequest, "No account name found")
	}

	var (
		acc *account.Account
		err error
	)

	acc = service.Account().LookupById(id)
	if acc == nil {
		acc, err = service.Account().UnsafeLookupById(id)
	}

	if err != nil {
		return common.HTTPError(c, http.StatusInternalServerError, "Failed to lookup account: "+err.Error())
	}

	if acc == nil {
		acc = account.Empty(id, "")
	}

	if acc.Name() != name {
		if acc.Name() != "" {
			service.Account().UpdateName(acc.Name(), name, acc.Id())
		} else {
			service.Account().Cache(acc)
		}

		acc.SetLastName(acc.Name())
		acc.SetName(name)

		go func() {
			err := service.Account().Update(acc)
			if err != nil {
				common.Log.Fatalf("Failed to update account: %s", err)
			}
		}()
	}

	acc.SetOnline(true)

	return c.Status(http.StatusOK).JSON(acc)
}

// Hook registers the route to the app
func Hook(app *fiber.App) {
	g := app.Group("/v1/account")
	g.Get("/:id/join/:name", AccountJoinRoute)
}
