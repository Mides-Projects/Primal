package common

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	RedisClient  *redis.Client = nil
	MongoClient  *mongo.Client = nil
	RedisChannel string
	APIKey       string
	Log          *log.Logger
)

func LoadRedis(redisUrl string) {
	if RedisClient != nil {
		Log.Panic("Redis already loaded")
	}

	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		Log.Panicf("Failed to parse Redis URL: %v", err)
	}

	RedisClient = redis.NewClient(opt)
	if _, err = RedisClient.Ping(context.Background()).Result(); err != nil {
		Log.Panicf("Failed to ping Redis: %v", err)
	}
}

func LoadMongo(uri string) {
	if MongoClient != nil {
		Log.Panic("MongoDB already loaded")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		Log.Panicf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		Log.Panicf("Failed to ping MongoDB: %v", err)
	}

	MongoClient = client
}
