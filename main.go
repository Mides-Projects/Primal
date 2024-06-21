package main

import (
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/common/config"
	"github.com/holypvp/primal/common/loader"
	"github.com/holypvp/primal/server"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"time"
)

func main() {
	file, err := os.ReadFile("config.yml")
	if err != nil {
		panic("config.yml not found")
	}

	configYaml := &config.Yaml{}
	err = yaml.Unmarshal(file, configYaml)
	if err != nil {
		panic("config.yml is invalid")
	}

	common.RedisChannel = configYaml.RedisChannel
	common.APIKey = configYaml.Key

	common.LoadMongo(configYaml.MongoUri)
	common.LoadRedis(configYaml.RedisUri)

	database := common.MongoClient.Database("api")

	server.Service().LoadGroups(database)
	server.Service().LoadServers(database)

	log.Println("App is running on port " + configYaml.Port + "...")

	loader.LoadAll(time.Now(), configYaml.Port)

	// route(router, "/players/{id}/lookup/{type}", playerRoute.LookupPlayer, "GET")
	// route(router, "/players/save", playerRoute.SavePlayer, "POST")

	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument

	// log.Fatal(http.ListenAndServe(":"+configYaml.Port, router))
}
