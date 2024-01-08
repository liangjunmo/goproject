package mutex

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/liangjunmo/goproject/internal/testutil"
)

func TestProvider(t *testing.T) {
	redisClient := testutil.InitRedis()
	defer redisClient.Close()

	sync := testutil.InitRedSync(redisClient)

	var mutexProvider *provider

	beforeTest := func(t *testing.T) {
		mutexProvider = newProvider(sync)
	}

	t.Run("ProvideCreateUserMutex", func(t *testing.T) {
		beforeTest(t)

		mutex := mutexProvider.ProvideCreateUserMutex(1)
		require.IsType(t, &createUserMutex{}, mutex)
	})
}
