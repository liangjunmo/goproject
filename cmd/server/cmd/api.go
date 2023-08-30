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
	"gorm.io/gorm"

	v1 "github.com/liangjunmo/goproject/internal/app/server/serverapi/v1"
	"github.com/liangjunmo/goproject/internal/app/server/serverconfig"
	"github.com/liangjunmo/goproject/internal/app/server/service/userservice"
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

		release := buildApi(router)
		defer release()

		server := &http.Server{
			Addr:    serverconfig.Config.Api.Addr,
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

func migrateDb(db *gorm.DB) {
	err := db.AutoMigrate(
		&userservice.User{},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func buildApi(router *gin.Engine) (release func()) {
	db := connectDb(serverconfig.Config.Debug)
	migrateDb(db)

	redisClient := connectRedis()

	release = func() {
		db, _ := db.DB()
		_ = db.Close()

		_ = redisClient.Close()
	}

	redisSync := newRedisSync(redisClient)

	userListService := userservice.NewListService(db)
	userReadService := userservice.NewReadService(db)
	userBusinessService := userservice.NewBusinessService(db, redisSync)
	userService := userservice.NewService(userListService, userReadService, userBusinessService)

	v1DefaultHandler := v1.NewDefaultHandler()
	v1AccountHandler := v1.NewAccountHandler(v1.NewAccountComponent(redisClient, userService))
	v1UserHandler := v1.NewUserHandler(userService)

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
