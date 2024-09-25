package grants

import (
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/model/grantsx"
	"github.com/holypvp/primal/service"
	"net/http"
)

func GrantsCreateRoute(c fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'name' parameter",
			"code":    http.StatusBadRequest,
		})
	}

	var body map[string]interface{}
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to bind request body",
			"code":    http.StatusBadRequest,
		})
	}

	g := &grantsx.Grant{}
	if err := g.Unmarshal(body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to unmarshal request body",
			"code":    http.StatusBadRequest,
		})
	}

	gaAdder := service.Grants().Lookup(g.AddedBy())
	if gaAdder == nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Who added the grant not found",
			"code":    http.StatusBadRequest,
		})
	}

	if gaAdder.Account().Id() != g.AddedBy() {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"message": "Who added the grant ID mismatch",
			"code":    http.StatusConflict,
		})
	}

	hgAdder := service.Grants().HighestGroupBy(gaAdder)
	if hgAdder == nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Who added the grant has no group",
			"code":    http.StatusBadRequest,
		})
	}

	// Retrieve the account of the player from our redis cache
	// but if they are online, we can fetch it from the RAM Cache
	acc, err := service.Account().UnsafeLookupByName(name)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"code":    http.StatusInternalServerError,
		})
	} else if acc == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Player not found",
			"code":    http.StatusNotFound,
		})
	}

	ga, err := service.Grants().UnsafeLookup(acc.Id())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"code":    http.StatusInternalServerError,
		})
	}

	if ga == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Player has no grantsx",
			"code":    http.StatusNotFound,
		})
	}

	if ga.Account().Id() != acc.Id() {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"message": "Player ID mismatch",
			"code":    http.StatusConflict,
		})
	}

	hg := service.Grants().HighestGroupBy(ga)
	if hg != nil && hg.Weight() > hgAdder.Weight() {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Who added the grant has lower group",
			"code":    http.StatusUnauthorized,
		})
	}

	go func() {
		if err := service.Grants().Save(acc.Id(), g); err != nil {
			common.Log.Fatalf("Failed to save grant: %v", err)
		}
	}()

	return c.SendStatus(http.StatusOK)
}
