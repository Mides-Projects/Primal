package grants

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/model/grantsx"
	"github.com/holypvp/primal/service"
	"net/http"
)

func CreateRoute(c fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing 'name' parameter",
		})
	}

	var body map[string]interface{}
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to bind request body",
		})
	}

	g := &grantsx.Grant{}
	if err := g.Unmarshal(body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to unmarshal request body",
		})
	}

	addrTrack := service.Grants().Lookup(g.AddedBy())
	if addrTrack == nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server: Our system could not find your trackable account",
		})
	}

	addrAcc := service.Account().LookupById(g.AddedBy())
	if addrAcc == nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server: Our system could not find your account",
		})
	}

	if addrAcc.Id() != g.AddedBy() {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"message": "Your ID doest not match with the Account ID",
		})
	}

	// Retrieve the account of the player from our redis cache
	// but if they are online, we can fetch it from the RAM Cache
	srcAcc, err := service.Account().UnsafeLookupByName(name)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server: " + err.Error(),
		})
	} else if srcAcc == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "We could not find the account you are looking for",
		})
	}

	srcTrack, err := service.Grants().UnsafeLookup(srcAcc.Id())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error: " + err.Error(),
		})
	}

	if srcTrack == nil {
		panic("The 'source' cannot be have nil tracker...")
	}

	if !addrAcc.Operator() {
		if srcAcc.Operator() {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": fmt.Sprintf("You cannot add a grant to %s because they are an operator", srcAcc.DisplayName()),
				"code":    http.StatusUnauthorized,
			})
		}

		hgAdder := service.Groups().LookupById(addrAcc.HighestGroup())
		if hgAdder == nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "An error occurred while trying to lookup your highest group",
				"code":    http.StatusBadRequest,
			})
		}

		hg := service.Groups().LookupById(srcAcc.HighestGroup())
		if hg != nil && hg.Weight() > hgAdder.Weight() {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": fmt.Sprintf("You cannot add a grant to %s because they have a higher group than you", srcAcc.DisplayName()),
				"code":    http.StatusUnauthorized,
			})
		}
	}

	go func() {
		if err = service.Grants().Save(srcAcc.Id(), g); err != nil {
			common.Log.Fatalf("Failed to save grant: %v", err)
		}
	}()

	return c.SendStatus(http.StatusOK)
}
