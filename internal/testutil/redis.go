package testutil

import (
	"context"
	"log"
	"os"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
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

func InitRedSync(redisClient *redis.Client) *redsync.Redsync {
	return redsync.New(goredis.NewPool(redisClient))
}
