package serverapi

import (
	golog "log"

	"github.com/gin-gonic/gin"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/liangjunmo/gotraceutil"

	"github.com/liangjunmo/goproject/internal/app/server"
	v1 "github.com/liangjunmo/goproject/internal/app/server/serverapi/v1"
	"github.com/liangjunmo/goproject/internal/app/server/serverconfig"
	"github.com/liangjunmo/goproject/internal/server/userservice"
)

func Build(router *gin.Engine) (release func()) {
	err := server.BuildTrace()
	if err != nil {
		golog.Fatal(err)
	}

	err = server.BuildLog()
	if err != nil {
		golog.Fatal(err)
	}

	db, err := server.BuildDb(serverconfig.Config.Debug)
	if err != nil {
		golog.Fatal(err)
	}

	err = db.AutoMigrate(
		&userservice.User{},
	)
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
	userService := userservice.NewUserService(db, redisSync)
	userHubService := userservice.NewHubService(userListService, userReadService, userService)

	v1DefaultHandler := v1.NewDefaultHandler()
	v1UserHandler := v1.NewUserHandler(userHubService)

	router.GET("/health", v1DefaultHandler.Health)

	router.Use(gotraceutil.GinMiddleware())

	router.GET("/api/v1/user/list", v1UserHandler.ListUser)
	router.GET("/api/v1/user/search", v1UserHandler.SearchUser)
	router.GET("/api/v1/user", v1UserHandler.GetUser)
	router.POST("/api/v1/user", v1UserHandler.CreateUser)

	return
}
