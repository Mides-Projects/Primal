package startup

import (
    "errors"
    "github.com/bytedance/sonic"
    "github.com/gofiber/fiber/v3"
    srvRoutes "github.com/holypvp/primal/server/routes"
    "github.com/holypvp/primal/service"
    "go.mongodb.org/mongo-driver/mongo"
)

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

    srvRoutes.Hook(app)

    return app.Listen(":3000")
}
