package router

import (
	"github.com/gin-gonic/gin"

	"mxshop_api/userop_web/api/message"
	"mxshop_api/userop_web/middlewares"
)

func InitMessageRouter(Router *gin.RouterGroup) {
	messageRouter := Router.Group("message").Use(middlewares.JWTAuth())
	{
		messageRouter.GET("",message.List)
		messageRouter.POST("",message.New)
	}
}
