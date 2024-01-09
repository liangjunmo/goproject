package mutex

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/liangjunmo/goproject/internal/testutil"
)

func TestMutexProvider(t *testing.T) {
	redisClient := testutil.InitRedis()
	defer redisClient.Close()

	sync := testutil.InitRedSync(redisClient)

	var provider *mutexProvider

	beforeTest := func(t *testing.T) {
		provider = newMutexProvider(sync)
	}

	t.Run("ProvideCreateUserMutex", func(t *testing.T) {
		beforeTest(t)

		mutex := provider.ProvideCreateUserMutex(1)
		require.IsType(t, &createUserMutex{}, mutex)
	})
}
