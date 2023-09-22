package userworker

import (
	"context"
	"sync"
	"time"

	"github.com/liangjunmo/gotraceutil"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/internal/service/userservice"
	"github.com/liangjunmo/goproject/internal/types"
)

type ListUserWorker struct {
	log         *log.Entry
	db          *gorm.DB
	userService userservice.Service
}

func NewListUserWorker(db *gorm.DB, userService userservice.Service) *ListUserWorker {
	return &ListUserWorker{
		log:         log.WithField("tag", "ListUserWorker"),
		db:          db,
		userService: userService,
	}
}

func (worker *ListUserWorker) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		worker.log.Info("quit")
		wg.Done()
	}()

	ticker := time.NewTicker(time.Second * 1)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			worker.handleTimer(gotraceutil.Trace(ctx))
		}
	}
}

func (worker *ListUserWorker) handleTimer(ctx context.Context) {
	var users []types.User

	err := worker.db.Model(&types.User{}).Find(&users).Error
	if err != nil {
		worker.log.WithContext(ctx).WithError(err).Error(err)
		return
	}

	for _, user := range users {
		worker.output(ctx, user.UID)
	}
}

func (worker *ListUserWorker) output(ctx context.Context, uid uint32) {
	user, err := worker.userService.GetUser(ctx, userservice.GetUserRequest{
		UID: uid,
	})
	if err != nil {
		worker.log.WithContext(ctx).WithError(err).Error(err)
		return
	}

	worker.log.WithContext(ctx).Infof("%+v", user)
}
