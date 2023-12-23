package api

import (
	"github.com/gin-gonic/gin"
	"github.com/liangjunmo/gotraceutil"

	"github.com/liangjunmo/goproject/internal/goproject/accountservice"
	"github.com/liangjunmo/goproject/internal/goproject/userservice"
)

func Router(config Config, engine *gin.Engine, accountService accountservice.Service, userService userservice.Service) {
	router(engine, newHandler(config, accountService, userService))
}

func router(router *gin.Engine, handler *handler) {
	router.GET("/ping", handler.Ping)

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
