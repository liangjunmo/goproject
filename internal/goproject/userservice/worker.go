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

type Worker interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
}

func ProvideWorker(db *gorm.DB, redisClient *redis.Client, userCenterClient usercenterproto.UserCenterClient) Worker {
	repository := newDefaultRepository(db)
	mutexProvider := newDefaultMutexProvider(redsync.New(goredis.NewPool(redisClient)))
	service := newDefaultService(repository, mutexProvider, userCenterClient)

	return newDefaultWorker(service)
}

type defaultWorker struct {
	service *defaultService
}

func newDefaultWorker(service *defaultService) *defaultWorker {
	return &defaultWorker{
		service: service,
	}
}

func (worker *defaultWorker) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go worker.runExample(ctx, wg)
}

func (worker *defaultWorker) runExample(ctx context.Context, wg *sync.WaitGroup) {
	log := logrus.WithField("tag", "goproject.userservice.worker.runExample")

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
			log.WithContext(gotraceutil.Trace(ctx)).Info("runExample")
		}
	}
}
