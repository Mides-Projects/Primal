package main

import (
	"github.com/holypvp/primal/common"
	"github.com/holypvp/primal/common/config"
	"github.com/holypvp/primal/common/startup"
	"gopkg.in/yaml.v2"
	"os"
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

	common.LoadMongo(conf.MongoUri)
	common.LoadRedis(conf.RedisUri)

	// Here I have a problem because startup depends of MongoDB and MongoDB depends on echo.Logger
	startup.LoadAll(time.Now(), conf.Port)

	startup.Shutdown()
}
