package bgroups

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/grantsx/model"
	"github.com/holypvp/primal/service"
	"net/http"
)

func create(c fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return common.HTTPError(c, http.StatusBadRequest, "No group name found")
	}

	if service.Groups().LookupByName(name) != nil {
		return common.HTTPError(c, http.StatusConflict, "Group already exists")
	}

	g := model.EmptyGroup(name)
	service.Groups().Cache(g)

	go func() {
		if err := service.Groups().Save(g); err != nil {
			common.Log.Fatalf("Failed to save group: %v", err)
		}
	}()

	return c.Status(http.StatusOK).JSON(g.Marshal())
}
