package usercenter

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/liangjunmo/gotraceutil"
	"google.golang.org/grpc"

	"github.com/liangjunmo/goproject/api/usercenterproto"
	"github.com/liangjunmo/goproject/internal/usercenter/mutex"
	"github.com/liangjunmo/goproject/internal/usercenter/repository"
	"github.com/liangjunmo/goproject/internal/usercenter/rpc"
	"github.com/liangjunmo/goproject/internal/usercenter/service"
)

type RPCServerConfig struct {
	Environment string
	Debug       bool
	Addr        string
	DB          DBConfig
	Redis       RedisConfig
}

func RunRPCServer(config RPCServerConfig) {
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

	userRepository := repository.NewUserRepository(db)

	mutexProvider := mutex.NewProvider(initRedSync(redisClient))

	userService := service.NewUserService(userRepository, mutexProvider)

	rpcServer := rpc.NewServer(userService)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(gotraceutil.GRPCUnaryServerInterceptor),
	)
	defer grpcServer.Stop()

	usercenterproto.RegisterUserCenterServer(grpcServer, rpcServer)

	listener, err := net.Listen("tcp", config.Addr)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := grpcServer.Serve(listener)
		if err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-c

	grpcServer.Stop()
}
