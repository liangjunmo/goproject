package testutil

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	// https://pkg.go.dev/github.com/go-redis/redis#example-ParseURL
	// GOPROJECT_TEST_REDIS="redis://user:password@localhost:6379/1"
	opts, err := redis.ParseURL(os.Getenv("GOPROJECT_TEST_REDIS"))
	if err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(opts)

	err = redisClient.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal(err)
	}

	return redisClient
}
