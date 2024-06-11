package main

import (
	"github.com/gorilla/mux"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/common/config"
	"github.com/holypvp/primal/common/loader"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
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

	common.LoadMongo(configYaml.MongoUri)
	common.LoadRedis(configYaml.RedisUri)

	router := mux.NewRouter().StrictSlash(true)

	loader.LoadAll(router)

	// route(router, "/players/{id}/lookup/{type}", playerRoute.LookupPlayer, "GET")
	// route(router, "/players/save", playerRoute.SavePlayer, "POST")

	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument

	log.Println("App is running on port " + configYaml.Port + "...")

	log.Fatal(http.ListenAndServe(":"+configYaml.Port, router))
}
