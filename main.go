package main

import (
	"context"
	"flag"
	"github.com/gorilla/mux"
	"github.com/holypvp/primal/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

func main() {
	uri := flag.String("uri", "null", "uri")
	port := flag.String("port", "3001", "port")

	flag.Parse()

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(*uri))
	if err != nil {
		panic("Failed to connect to MongoDB: " + err.Error())
	}

	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		panic("Failed to ping MongoDB: " + err.Error())
	}

	router := mux.NewRouter().StrictSlash(true)

	server.LoadAll(router)

	// route(router, "/players/{id}/lookup/{type}", playerRoute.LookupPlayer, "GET")
	// route(router, "/players/save", playerRoute.SavePlayer, "POST")

	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	log.Fatal(http.ListenAndServe(":"+*port, router))
}
