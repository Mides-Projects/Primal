package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/common/loader"
	"log"
	"net/http"
)

func main() {
	uri := *flag.String("uri", "null", "uri")
	redisUrl := *flag.String("redis", "null", "redis")
	port := *flag.String("port", "3001", "port")

	flag.Parse()

	if uri == "null" {
		panic("uri is required")
	}

	if redisUrl == "null" {
		panic("redis is required")
	}

	common.LoadMongo(uri)
	common.LoadRedis(redisUrl)

	router := mux.NewRouter().StrictSlash(true)

	loader.LoadAll(router)

	// route(router, "/players/{id}/lookup/{type}", playerRoute.LookupPlayer, "GET")
	// route(router, "/players/save", playerRoute.SavePlayer, "POST")

	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument

	log.Println("App is running on port " + port + "...")

	log.Fatal(http.ListenAndServe(":"+port, router))
}
