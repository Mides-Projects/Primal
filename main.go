package main

import (
    "github.com/holypvp/primal/common"
    "github.com/holypvp/primal/common/config"
    "github.com/holypvp/primal/common/startup"
    "gopkg.in/yaml.v2"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    file, err := os.ReadFile("config.yml")
    if err != nil {
        panic("config.yml not found")
    }

    conf := &config.Yaml{}
    err = yaml.Unmarshal(file, conf)
    if err != nil {
        panic("config.yml is invalid")
    }

    common.RedisChannel = conf.RedisChannel
    common.APIKey = conf.Key

    // Here I have a problem because startup depends of MongoDB and MongoDB depends on echo.Logger
    startup.LoadAll(time.Now(), conf.Port)

    common.LoadMongo(conf.MongoUri)
    common.LoadRedis(conf.RedisUri)

    // server.Service().LoadGroups(database)
    // server.Service().LoadServers(database)

    log.Println("App is running on port " + conf.Port + "...")

    go func() {
        sig := make(chan os.Signal, 1)
        signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
        <-sig
        log.Println("Shutting down...")
        startup.Shutdown()
    }()

    // route(router, "/players/{id}/lookup/{type}", playerRoute.LookupPlayer, "GET")
    // route(router, "/players/save", playerRoute.SavePlayer, "POST")

    // finally, instead of passing in nil, we want
    // to pass in our newly created router as the second
    // argument

    // log.Fatal(http.ListenAndServe(":"+conf.Port, router))
}
