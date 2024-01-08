package manager

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

type RedisManager interface {
	GetLoginFailedCount(ctx context.Context, username string) (uint32, error)
	SetLoginFailedCount(ctx context.Context, username string) error
	DelLoginFailedCount(ctx context.Context, username string) error
	GetLoginTicket(ctx context.Context, ticket string) (uid uint32, exist bool, err error)
	SetLoginTicket(ctx context.Context, ticket string, uid uint32, expiration time.Duration) error
}

func NewRedisManager(redisClient *redis.Client) RedisManager {
	return newRedisManager(redisClient)
}

type redisManager struct {
	redisClient *redis.Client
}

func newRedisManager(redisClient *redis.Client) *redisManager {
	return &redisManager{
		redisClient: redisClient,
	}
}

func (manager *redisManager) GetLoginFailedCount(ctx context.Context, username string) (uint32, error) {
	val, err := manager.redisClient.Get(ctx, manager.keyLoginFailedCount(username)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, err
	}

	return cast.ToUint32(val), nil
}

func (manager *redisManager) SetLoginFailedCount(ctx context.Context, username string) error {
	_, err := manager.redisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Incr(ctx, manager.keyLoginFailedCount(username))
		pipe.Expire(ctx, manager.keyLoginFailedCount(username), time.Minute*5)

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (manager *redisManager) DelLoginFailedCount(ctx context.Context, username string) error {
	err := manager.redisClient.Del(ctx, manager.keyLoginFailedCount(username)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (manager *redisManager) GetLoginTicket(ctx context.Context, ticket string) (uid uint32, exist bool, err error) {
	val, err := manager.redisClient.Get(ctx, manager.keyLoginTicket(ticket)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, false, err
	}

	if errors.Is(err, redis.Nil) {
		return 0, false, nil
	}

	return cast.ToUint32(val), true, nil
}

func (manager *redisManager) SetLoginTicket(ctx context.Context, ticket string, uid uint32, expiration time.Duration) error {
	err := manager.redisClient.Set(ctx, manager.keyLoginTicket(ticket), uid, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (manager *redisManager) keyLoginFailedCount(username string) string {
	return fmt.Sprintf("goproject-login-failed-count-%s", username)
}

func (manager *redisManager) keyLoginTicket(ticket string) string {
	return fmt.Sprintf("goproject-login-ticket-%s", ticket)
}
