package goproject

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liangjunmo/gotraceutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/liangjunmo/goproject/api/usercenterproto"
	"github.com/liangjunmo/goproject/internal/goproject/api"
	"github.com/liangjunmo/goproject/internal/goproject/manager"
	"github.com/liangjunmo/goproject/internal/goproject/mutex"
	"github.com/liangjunmo/goproject/internal/goproject/repository"
	"github.com/liangjunmo/goproject/internal/goproject/service"
	"github.com/liangjunmo/goproject/internal/goproject/usecase"
)

type APIServerConfig struct {
	Environment             string
	Debug                   bool
	JWTKey                  string
	Addr                    string
	UserCenterRPCServerAddr string
	DB                      DBConfig
	Redis                   RedisConfig
}

func RunAPIServer(config APIServerConfig) {
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

	userCenterConn, err := grpc.Dial(
		config.UserCenterRPCServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(gotraceutil.GRPCUnaryClientInterceptor),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer userCenterConn.Close()

	userRepository := repository.NewUserRepository(db)

	mutexProvider := mutex.NewProvider(initRedSync(redisClient))

	redisManager := manager.NewRedisManager(redisClient)

	userCenterClient := usercenterproto.NewUserCenterClient(userCenterConn)

	userService := service.NewUserService(userRepository, mutexProvider, userCenterClient)

	accountUseCase := usecase.NewAccountUseCase(
		usecase.AccountUseCaseConfig{JWTKey: config.JWTKey},
		redisManager,
		userService,
	)

	engine := gin.Default()

	api.Router(
		api.Config{
			Debug:        config.Debug,
			TracingIDKey: tracingKeys[0],
		},
		engine,
		accountUseCase,
		userService,
	)

	server := &http.Server{
		Addr:    config.Addr,
		Handler: engine,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
