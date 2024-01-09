package goproject

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/liangjunmo/goproject/api/usercenterproto"
	"github.com/liangjunmo/goproject/internal/goproject/mutex"
	"github.com/liangjunmo/goproject/internal/goproject/repository"
	"github.com/liangjunmo/goproject/internal/goproject/service"
)

type WorkerServerConfig struct {
	Environment             string
	Debug                   bool
	UserCenterRPCServerAddr string
	DB                      DBConfig
	Redis                   RedisConfig
}

func RunWorkerServer(config WorkerServerConfig) {
	tracingKeys := []string{"TracingID"}

	initTracing(tracingKeys)
	initLogger()

	db := initDB(config.DB, config.Debug)
	defer func() {
		db, _ := db.DB()
		_ = db.Close()
	}()

	redisClient := initRedis(config.Redis)
	defer redisClient.Close()

	userCenterConn, err := grpc.Dial(config.UserCenterRPCServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer userCenterConn.Close()

	mutexProvider := mutex.NewMutexProvider(initRedSync(redisClient))

	userRepository := repository.NewUserRepository(db)

	userCenterClient := usercenterproto.NewUserCenterClient(userCenterConn)

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	service.RunUserScheduler(ctx, wg, mutexProvider, userRepository, userCenterClient)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-c

	cancel()
	wg.Wait()
}
