package mutex

import (
	"context"
	"fmt"

	"github.com/go-redsync/redsync/v4"
)

type Mutex interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) (ok bool, err error)
}

type createUserMutex struct {
	mutex *redsync.Mutex
}

func newCreateUserMutex(sync *redsync.Redsync, username string) *createUserMutex {
	mutex := &createUserMutex{}

	mutex.mutex = sync.NewMutex(
		mutex.key(username),
		redsync.WithTries(1),
	)

	return mutex
}

func (mutex *createUserMutex) Lock(ctx context.Context) error {
	return mutex.mutex.LockContext(ctx)
}

func (mutex *createUserMutex) Unlock(ctx context.Context) (ok bool, err error) {
	return mutex.mutex.UnlockContext(ctx)
}

func (mutex *createUserMutex) key(username string) string {
	return fmt.Sprintf("goproject-usercenter-create-user-mutex-%s", username)
}
