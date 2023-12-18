package userservice

import (
	"context"
	"testing"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/stretchr/testify/require"

	"github.com/liangjunmo/goproject/internal/goproject/testutil"
)

func TestDefaultMutexProvider(t *testing.T) {
	redisClient := testutil.InitRedis()
	sync := redsync.New(goredis.NewPool(redisClient))

	t.Run("ProvideCreateUserMutex", func(t *testing.T) {
		mutexProvider := newDefaultMutexProvider(sync)

		mutex := mutexProvider.ProvideCreateUserMutex(1)
		require.IsType(t, &createUserMutex{}, mutex)
	})

	redisClient.Close()
}

func TestCreateUserMutex(t *testing.T) {
	redisClient := testutil.InitRedis()
	sync := redsync.New(goredis.NewPool(redisClient))

	t.Run("LockAndUnlock", func(t *testing.T) {
		mutex := newCreateUserMutex(sync, 1)

		ctx := context.Background()

		err := mutex.Lock(ctx)
		require.Nil(t, err)

		ok, err := mutex.Unlock(ctx)
		require.Nil(t, err)
		require.True(t, ok)
	})

	redisClient.Close()
}
