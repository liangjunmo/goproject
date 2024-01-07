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
	"github.com/liangjunmo/goproject/internal/usercenter/rpc"
	"github.com/liangjunmo/goproject/internal/usercenter/service/userservice/userserviceimpl"
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
	redisClient := initRedis(config.Redis)

	userService := userserviceimpl.ProvideService(db, redisClient)

	rpcServer := rpc.NewServer(userService)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(gotraceutil.GRPCUnaryServerInterceptor),
	)

	usercenterproto.RegisterUserCenterServer(grpcServer, rpcServer)

	defer func() {
		grpcServer.Stop()

		db, _ := db.DB()
		_ = db.Close()

		_ = redisClient.Close()
	}()

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
