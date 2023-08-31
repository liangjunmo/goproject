package v1

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"

	"github.com/liangjunmo/goproject/internal/app/rediskey"
	"github.com/liangjunmo/goproject/internal/app/server/codes"
)

func RedisGetLoginFailedCount(ctx context.Context, redisClient *redis.Client, username string) (uint32, error) {
	val, err := redisClient.Get(ctx, rediskey.LoginFailedCount(username)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return cast.ToUint32(val), nil
}

func RedisSetLoginFailedCount(ctx context.Context, redisClient *redis.Client, username string) error {
	_, err := redisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Incr(ctx, rediskey.LoginFailedCount(username))
		pipe.Expire(ctx, rediskey.LoginFailedCount(username), time.Minute*5)

		return nil
	})
	if err != nil {
		return fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return nil
}

func RedisDelLoginFailedCount(ctx context.Context, redisClient *redis.Client, username string) error {
	err := redisClient.Del(ctx, rediskey.LoginFailedCount(username)).Err()
	if err != nil {
		return fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return nil
}

func RedisGetLoginTicket(ctx context.Context, redisClient *redis.Client, ticket string) (uint32, bool, error) {
	val, err := redisClient.Get(ctx, rediskey.LoginTicket(ticket)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, false, fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	if errors.Is(err, redis.Nil) {
		return 0, false, nil
	}

	return cast.ToUint32(val), true, nil
}

func RedisSetLoginTicket(ctx context.Context, redisClient *redis.Client, ticket string, uid uint32, expiration time.Duration) error {
	err := redisClient.Set(ctx, rediskey.LoginTicket(ticket), uid, expiration).Err()
	if err != nil {
		return fmt.Errorf("%w: %v", codes.InternalServerError, err)
	}

	return nil
}
