package common

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	RedisClient  *redis.Client = nil
	MongoClient  *mongo.Client = nil
	RedisChannel string
	APIKey       string
)

func LoadRedis(redisUrl string) {
	if RedisClient != nil {
		panic("Redis already loaded")
	}

	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic("Failed to parse Redis URL: " + err.Error())
	}

	RedisClient = redis.NewClient(opt)
	_, err = RedisClient.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to ping Redis: " + err.Error())
	}
}

func LoadMongo(uri string) {
	if MongoClient != nil {
		panic("Mongo already loaded")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic("Failed to connect to MongoDB: " + err.Error())
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic("Failed to ping MongoDB: " + err.Error())
	}

	MongoClient = client
}

func WrapPayload(pid string, payload interface{}) ([]byte, error) {
	return json.Marshal(NewPayload(pid, 0, payload))
}
