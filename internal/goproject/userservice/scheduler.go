package userservice

import (
	"context"
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/liangjunmo/gotraceutil"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/api/usercenterproto"
)

func RunScheduler(ctx context.Context, wg *sync.WaitGroup, db *gorm.DB, redisClient *redis.Client, userCenterClient usercenterproto.UserCenterClient) {
	service := newDefaultService(
		newDefaultRepository(db),
		newDefaultMutexProvider(redsync.New(goredis.NewPool(redisClient))),
		userCenterClient,
	)

	runScheduler(ctx, wg, service)
}

func runScheduler(ctx context.Context, wg *sync.WaitGroup, service *defaultService) {
	wg.Add(1)
	go jobToRunExample(ctx, wg, service)
}

func jobToRunExample(ctx context.Context, wg *sync.WaitGroup, service *defaultService) {
	log := logrus.WithField("tag", "goproject.userservice.scheduler.jobToRunExample")

	defer func() {
		log.Info("quit")
		wg.Done()
	}()

	ticker := time.NewTicker(time.Second * 1)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			service.taskToRunExample(gotraceutil.Trace(ctx), log)
		}
	}
}
