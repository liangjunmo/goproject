package redismutex

import (
	"github.com/go-redsync/redsync/v4"

	"github.com/liangjunmo/goproject/internal/rediskey"
)

func NewCreateUserCenterUserMutex(redisSync *redsync.Redsync, username string) *redsync.Mutex {
	return redisSync.NewMutex(
		rediskey.MutexCreateUserCenterUser(username),
		redsync.WithTries(1),
	)
}

func NewCreateUserMutex(redisSync *redsync.Redsync, uid uint32) *redsync.Mutex {
	return redisSync.NewMutex(
		rediskey.MutexCreateUser(uid),
		redsync.WithTries(1),
	)
}
