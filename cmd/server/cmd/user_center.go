package cmd

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/liangjunmo/goproject/api/usercenterproto"
	"github.com/liangjunmo/goproject/internal/app/server/config"
	"github.com/liangjunmo/goproject/internal/app/server/usercenter"
	"github.com/liangjunmo/goproject/internal/service/usercenterservice"
)

func init() {
	rootCmd.AddCommand(userCenterCmd)
}

var userCenterCmd = &cobra.Command{
	Use:   "user-center",
	Short: "goproject-server user center cli tool",
	Long:  "goproject-server user center cli tool",
	Run: func(cmd *cobra.Command, args []string) {
		server, release := buildUserCenter()
		defer release()

		lis, err := net.Listen("tcp", config.Config.UserCenter.Addr)
		if err != nil {
			log.Fatal(err)
		}

		s := grpc.NewServer()
		usercenterproto.RegisterUserCenterServer(s, server)

		go func() {
			err := s.Serve(lis)
			if err != nil {
				log.Fatal(err)
			}
		}()

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		<-c

		s.Stop()
	},
}

func buildUserCenter() (server *usercenter.Server, release func()) {
	db := connectDB(config.Config.Debug)
	migrateDB(db)

	redisClient := connectRedis()

	release = func() {
		db, _ := db.DB()
		_ = db.Close()

		_ = redisClient.Close()
	}

	redisSync := newRedisSync(redisClient)

	userCenterService := usercenterservice.NewService(db, redisSync)
	server = usercenter.NewServer(userCenterService)

	return server, release
}
