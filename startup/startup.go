package startup

import (
	"errors"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/keyauth"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/holypvp/primal/common"
	grantRoutes "github.com/holypvp/primal/grantsx/routes"
	srvRoutes "github.com/holypvp/primal/server/routes"
	"github.com/holypvp/primal/service"
	"go.mongodb.org/mongo-driver/mongo"
)

// Hook is the entry point for the Primal API
func Hook(db *mongo.Database) error {
	if err := service.LoadServers(db); err != nil {
		return errors.Join(errors.New("failed to load servers"), err)
	} else if err := service.LoadGroups(db); err != nil {
		return errors.Join(errors.New("failed to load server groups"), err)
	} else if err := service.Groups().Hook(db); err != nil {
		return errors.Join(errors.New("failed to hook 'BungeeGroups'"), err)
	} else if err := service.Grants().Hook(db); err != nil {
		return errors.Join(errors.New("failed to hook 'GrantsX'"), err)
	}

	app := fiber.New(fiber.Config{
		JSONEncoder:   sonic.Marshal,
		JSONDecoder:   sonic.Unmarshal,
		StrictRouting: true,
		ServerHeader:  "Primal v0.0.1",
		AppName:       "Primal v0.0.1",
	})

	app.Use(recover.New())
	app.Use(keyauth.New(keyauth.Config{
		KeyLookup: "header:x-api-key",
		SuccessHandler: func(c fiber.Ctx) error {
			return c.Next()
		},
		ErrorHandler: func(c fiber.Ctx, err error) error {
			if errors.Is(err, keyauth.ErrMissingOrMalformedAPIKey) {
				return c.Status(fiber.StatusUnauthorized).SendString("missing or malformed API key")
			} else {
				return c.Status(fiber.StatusUnauthorized).SendString("invalid or expired API key")
			}
		},
		Validator: func(c fiber.Ctx, input string) (bool, error) {
			if common.APIKey == input {
				return true, nil
			} else if input == "" {
				return false, errors.New("missing or malformed API key")
			} else {
				return false, errors.New("invalid or expired API key")
			}
		},
	}))

	grantRoutes.Hook(app)
	srvRoutes.Hook(app)

	return app.Listen(":3000")
}
