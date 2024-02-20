package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
)

func InitRedis() error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()
	// Test the connection
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	fmt.Println("connected to redis successfully")
	return nil
}

func GetClient() (*redis.Client, error) {
	if redisClient == nil {
		return nil, fmt.Errorf("error in fetching the redis client")
	}
	return redisClient, nil
}
