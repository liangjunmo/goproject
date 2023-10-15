package cmd

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
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"

	"github.com/liangjunmo/goproject/api/usercenterproto"
	"github.com/liangjunmo/goproject/internal/manager"
	v1 "github.com/liangjunmo/goproject/internal/server/api/v1"
	"github.com/liangjunmo/goproject/internal/server/config"
	"github.com/liangjunmo/goproject/internal/server/usercenter"
	"github.com/liangjunmo/goproject/internal/service/userservice"
	"github.com/liangjunmo/goproject/internal/types"
)

func init() {
	rootCmd.AddCommand(apiCmd)
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "goproject-server api cli tool",
	Long:  "goproject-server api cli tool",
	Run: func(cmd *cobra.Command, args []string) {
		router := gin.Default()

		release := buildAPI(router)
		defer release()

		server := &http.Server{
			Addr:    config.Config.API.Addr,
			Handler: router,
		}

		go func() {
			err := server.ListenAndServe()
			if err == http.ErrServerClosed {
				log.Println("http server closed")
			} else {
				log.Fatal(err)
			}
		}()

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		<-c

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func migrateDB(db *gorm.DB) {
	err := db.AutoMigrate(
		&types.UserCenterUser{},
		&types.User{},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func buildAPI(router *gin.Engine) (release func()) {
	db := connectDB(config.Config.Debug)
	migrateDB(db)

	redisClient := connectRedis()

	userCenterConn, err := grpc.Dial(config.Config.API.UserCenterAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	release = func() {
		db, _ := db.DB()
		_ = db.Close()

		_ = redisClient.Close()

		userCenterConn.Close()
	}

	redisSync := newRedisSync(redisClient)

	userCenterClient := usercenterproto.NewUserCenterClient(userCenterConn)
	userCenterServiceAdapter := usercenter.NewAdapter(userCenterClient)

	userService := userservice.NewService(db, redisSync)
	userManager := manager.NewUserManager(userCenterServiceAdapter, userService)

	v1DefaultHandler := v1.NewDefaultHandler()
	v1AccountHandler := v1.NewAccountHandler(v1.NewAccountComponent(redisClient, userManager))
	v1UserHandler := v1.NewUserHandler(userManager)

	router.GET("/health", v1DefaultHandler.Health)

	router.Use(gotraceutil.GinMiddleware())

	router.POST("/api/v1/login", v1AccountHandler.Login)
	router.POST("/api/v1/token", v1AccountHandler.CreateToken)

	v1AuthGroup := router.Group("", v1AccountHandler.AuthMiddleware)
	{
		v1AuthGroup.GET("/api/v1/user/list", v1UserHandler.ListUser)
		v1AuthGroup.GET("/api/v1/user/search", v1UserHandler.SearchUser)
		v1AuthGroup.GET("/api/v1/user", v1UserHandler.GetUser)
		v1AuthGroup.POST("/api/v1/user", v1UserHandler.CreateUser)
	}

	return
}
