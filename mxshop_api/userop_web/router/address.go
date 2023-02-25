package router

import (
	"github.com/gin-gonic/gin"

	"mxshop_api/userop_web/api/address"
	"mxshop_api/userop_web/middlewares"
)

func InitAddressRouter(Router *gin.RouterGroup) {
	addressRouter := Router.Group("address").Use(middlewares.JWTAuth())
	{
		addressRouter.GET("", address.List)
		addressRouter.POST("", address.New)
		addressRouter.PUT("/:id", address.Update)
		addressRouter.DELETE("/:id", address.Delete)
	}
}
