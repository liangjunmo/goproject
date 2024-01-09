package mutex

import (
	"github.com/go-redsync/redsync/v4"
)

type MutexProvider interface {
	ProvideCreateUserMutex(username string) Mutex
}

func NewMutexProvider(sync *redsync.Redsync) MutexProvider {
	return newMutexProvider(sync)
}

type mutexProvider struct {
	sync *redsync.Redsync
}

func newMutexProvider(sync *redsync.Redsync) *mutexProvider {
	return &mutexProvider{
		sync: sync,
	}
}

func (provider *mutexProvider) ProvideCreateUserMutex(username string) Mutex {
	return newCreateUserMutex(provider.sync, username)
}
