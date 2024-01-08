package mutex

import (
	"github.com/go-redsync/redsync/v4"
)

type Provider interface {
	ProvideCreateUserMutex(uid uint32) Mutex
}

func NewProvider(sync *redsync.Redsync) Provider {
	return newProvider(sync)
}

type provider struct {
	sync *redsync.Redsync
}

func newProvider(sync *redsync.Redsync) *provider {
	return &provider{
		sync: sync,
	}
}

func (provider *provider) ProvideCreateUserMutex(uid uint32) Mutex {
	return newCreateUserMutex(provider.sync, uid)
}
