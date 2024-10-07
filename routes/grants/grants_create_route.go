package grants

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/helper"
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

	if service.Grants().Lookup(g.AddedBy()) == nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server: Our system could not find your trackable account",
		})
	}

	// SPI = Source Player Info
	spi := service.Player().LookupById(g.AddedBy())
	if spi == nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server: Our system could not find your account",
		})
	}

	if spi.Id() != g.AddedBy() {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"message": "Your ID doest not match with the Player ID",
		})
	}

	// Retrieve the account of the player from our redis cache
	// but if they are online, we can fetch it from the RAM Cache
	// TPI = Target Player Info
	tpi, err := service.Player().UnsafeLookupByName(name, false)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server: " + err.Error(),
		})
	} else if tpi == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "We could not find the player you are looking for",
		})
	}

	if track, err := service.Grants().UnsafeLookup(tpi.Id(), false); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server: " + err.Error(),
		})
	} else if track == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "We could not find the player you are looking for",
		})
	}

	if !spi.Operator() {
		if tpi.Operator() {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": fmt.Sprintf("You cannot add a grant to %s because they are an operator", tpi.DisplayName()),
				"code":    http.StatusUnauthorized,
			})
		} else if ok, err := helper.HighestThan(spi.HighestGroup(), tpi.HighestGroup()); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error: " + err.Error(),
			})
		} else if !ok {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": fmt.Sprintf("You cannot add a grant to %s because they have a higher group than you", tpi.DisplayName()),
				"code":    http.StatusUnauthorized,
			})
		}
	}

	go func() {
		if err = service.Grants().Save(tpi.Id(), g); err != nil {
			common.Log.Fatalf("Failed to save grant: %v", err)
		}
	}()

	return c.SendStatus(http.StatusOK)
}
