package main

import (
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/common/config"
	"github.com/holypvp/primal/startup"
	"gopkg.in/yaml.v2"
	"log"
	"os"
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

	common.Log = log.New(os.Stdout, "Primal", log.LstdFlags)

	common.RedisChannel = conf.RedisChannel
	common.APIKey = conf.Key

	common.LoadMongo(conf.MongoUri)
	common.LoadRedis(conf.RedisUri)

	db := common.MongoClient.Database("api")

	log.Print("Running!")
	log.Print(startup.Hook(db))
}
