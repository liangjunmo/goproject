package redismutex

import (
	"github.com/go-redsync/redsync/v4"

	"github.com/liangjunmo/goproject/internal/app/server/rediskey"
)

func NewMutexCreateUser(redisSync *redsync.Redsync, username string) *redsync.Mutex {
	return redisSync.NewMutex(
		rediskey.MutexCreateUser(username),
		redsync.WithTries(1),
	)
}