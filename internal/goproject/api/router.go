package api

import (
	"github.com/gin-gonic/gin"
	"github.com/liangjunmo/gotraceutil"
)

func Router(router *gin.Engine, handler *Handler) {
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
