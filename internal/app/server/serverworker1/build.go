package serverworker1

import (
	"context"
	golog "log"
	"sync"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"

	"github.com/liangjunmo/goproject/internal/app/server"
	"github.com/liangjunmo/goproject/internal/app/server/service/userservice"
	"github.com/liangjunmo/goproject/internal/app/server/worker/userworker"
)

func Build(ctx context.Context, wg *sync.WaitGroup) (release func()) {
	err := server.BuildTrace()
	if err != nil {
		golog.Fatal(err)
	}

	err = server.BuildLog()
	if err != nil {
		golog.Fatal(err)
	}

	db, err := server.BuildDb(true)
	if err != nil {
		golog.Fatal(err)
	}

	redisClient, err := server.BuildRedis()
	if err != nil {
		golog.Fatal(err)
	}

	release = func() {
		db, _ := db.DB()
		_ = db.Close()

		_ = redisClient.Close()
	}

	redisSync := redsync.New(goredis.NewPool(redisClient))

	userListService := userservice.NewListService(db)
	userReadService := userservice.NewReadService(db)
	userBusinessService := userservice.NewBusinessService(db, redisSync)
	userService := userservice.NewService(userListService, userReadService, userBusinessService)

	wg.Add(1)
	go userworker.NewListUserWorker(db, userService).Run(ctx, wg)

	return
}
