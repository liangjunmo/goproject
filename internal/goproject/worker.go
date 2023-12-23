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
	"github.com/liangjunmo/goproject/internal/goproject/userservice"
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
	redisClient := initRedis(config.Redis)

	userCenterConn, err := grpc.Dial(config.UserCenterRPCServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		userCenterConn.Close()

		db, _ := db.DB()
		_ = db.Close()

		_ = redisClient.Close()
	}()

	userCenterClient := usercenterproto.NewUserCenterClient(userCenterConn)

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	userservice.RunScheduler(ctx, wg, db, redisClient, userCenterClient)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-c

	cancel()
	wg.Wait()
}
