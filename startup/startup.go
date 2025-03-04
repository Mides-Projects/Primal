package startup

import (
    "errors"
    "github.com/bytedance/sonic"
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/keyauth"
    "github.com/gofiber/fiber/v3/middleware/logger"
    "github.com/gofiber/fiber/v3/middleware/recover"
    "github.com/holypvp/primal/common"
    accRoutes "github.com/holypvp/primal/routes/account"
    grantRoutes "github.com/holypvp/primal/routes/bgroups"
    srvRoutes "github.com/holypvp/primal/routes/server"
    "github.com/holypvp/primal/service"
    "go.mongodb.org/mongo-driver/mongo"
    "log"
    "time"
)

// Hook is the entry point for the Primal API
func Hook(db *mongo.Database) error {
    if err := service.LoadServers(db); err != nil {
        return errors.Join(errors.New("failed to load servers"), err)
    } else if err := service.LoadGroups(db); err != nil {
        return errors.Join(errors.New("failed to load server bgroups"), err)
    } else if err := service.Groups().Hook(db); err != nil {
        return errors.Join(errors.New("failed to hook 'BungeeGroups'"), err)
    } else if err := service.Grants().Hook(db); err != nil {
        return errors.Join(errors.New("failed to hook 'GrantsX'"), err)
    } else if err := service.Account().Hook(db); err != nil {
        return errors.Join(errors.New("failed to hook 'AccountService'"), err)
    }

    app := fiber.New(fiber.Config{
        JSONEncoder:   sonic.Marshal,
        JSONDecoder:   sonic.Unmarshal,
        StrictRouting: true,
        ServerHeader:  "Primal v0.0.1",
        AppName:       "Primal v0.0.1",
    })

    app.Use(recover.New())
    app.Use(logger.New())
    app.Use(keyauth.New(keyauth.Config{
        KeyLookup: "header:x-api-key",
        SuccessHandler: func(c fiber.Ctx) error {
            return c.Next()
        },
        ErrorHandler: func(c fiber.Ctx, err error) error {
            if errors.Is(err, keyauth.ErrMissingOrMalformedAPIKey) {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                    "message": "missing or malformed API key",
                    "code":    fiber.StatusUnauthorized,
                })
            } else {
                return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                    "message": "invalid or expired API key",
                    "code":    fiber.StatusForbidden,
                })
            }
        },
        Validator: func(_ fiber.Ctx, input string) (bool, error) {
            if common.APIKey == input {
                return true, nil
            } else if input == "" {
                return false, errors.New("missing or malformed API key")
            } else {
                return false, errors.New("invalid or expired API key")
            }
        },
    }))

    log.Print("Primal API is now running")

    grantRoutes.Hook(app)
    srvRoutes.Hook(app)
    accRoutes.Hook(app)

    // I do this outside of the goroutine because I want stop the ticker when the server stops.
    tc := time.NewTicker(time.Millisecond * 50)
    defer tc.Stop()

    go func() {
        for {
            for range tc.C {
                service.Account().DoTTLTick()
            }
        }
    }()

    defer func(app *fiber.App) {
        err := app.Shutdown()
        if err != nil {
            log.Fatalf("Failed to shutdown Primal API: %v", err)
        } else {
            log.Print("Primal API has been shutdown")
        }
    }(app)

    return app.Listen(common.Port)
}
