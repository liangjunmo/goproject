package userserviceimpl

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/liangjunmo/goproject/internal/testutil"
)

func TestDefaultMutexProvider(t *testing.T) {
	redisClient := testutil.InitRedis()
	defer redisClient.Close()

	sync := testutil.InitRedSync(redisClient)

	var mutexProvider *defaultMutexProvider

	beforeTest := func(t *testing.T) {
		mutexProvider = newDefaultMutexProvider(sync)
	}

	t.Run("ProvideCreateUserMutex", func(t *testing.T) {
		beforeTest(t)

		mutex := mutexProvider.ProvideCreateUserMutex(1)
		require.IsType(t, &createUserMutex{}, mutex)
	})

}

func TestCreateUserMutex(t *testing.T) {
	redisClient := testutil.InitRedis()
	defer redisClient.Close()

	sync := testutil.InitRedSync(redisClient)

	var (
		mutex *createUserMutex
		ctx   context.Context
	)

	beforeTest := func(t *testing.T) {
		mutex = newCreateUserMutex(sync, 1)

		ctx = context.Background()
	}

	t.Run("LockAndUnlock", func(t *testing.T) {
		beforeTest(t)

		err := mutex.Lock(ctx)
		require.Nil(t, err)

		ok, err := mutex.Unlock(ctx)
		require.Nil(t, err)
		require.True(t, ok)
	})
}
