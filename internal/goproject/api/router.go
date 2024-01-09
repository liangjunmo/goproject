package api

import (
	"github.com/gin-gonic/gin"
	"github.com/liangjunmo/gotraceutil"

	swaggerFiles "github.com/swaggo/files"
	swaggerGin "github.com/swaggo/gin-swagger"

	"github.com/liangjunmo/goproject/internal/goproject/service"
	"github.com/liangjunmo/goproject/internal/goproject/usecase"

	_ "github.com/liangjunmo/goproject/internal/goproject/api/swagger"
)

//	@title	GoProject API

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

func Router(config Config, engine *gin.Engine, accountUseCase usecase.AccountUseCase, userService service.UserService) {
	handler := newHandler(config, accountUseCase, userService)

	router(config.Debug, engine, handler)
}

func router(debug bool, router *gin.Engine, handler *handler) {
	router.GET("/ping", handler.Ping)

	if debug {
		router.GET("/swagger/*any", swaggerGin.WrapHandler(swaggerFiles.Handler))
	}

	router.Use(gotraceutil.GinMiddleware())

	router.POST("/api/v1/login", handler.Login)
	router.POST("/api/v1/token", handler.CreateToken)

	{
		router := router.Group("", handler.Authorize)
		{
			router.GET("/api/v1/user/list", handler.ListUser)
			router.GET("/api/v1/user/search", handler.SearchUser)
			router.GET("/api/v1/user/:uid", handler.GetUser)
			router.POST("/api/v1/user", handler.CreateUser)
		}
	}
}
