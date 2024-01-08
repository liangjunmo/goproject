package mutex

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/liangjunmo/goproject/internal/testutil"
)

func TestCreateUserMutex(t *testing.T) {
	redisClient := testutil.InitRedis()
	defer redisClient.Close()

	sync := testutil.InitRedSync(redisClient)

	var (
		mutex *createUserMutex
		ctx   context.Context
	)

	beforeTest := func(t *testing.T) {
		mutex = newCreateUserMutex(sync, "user")

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
