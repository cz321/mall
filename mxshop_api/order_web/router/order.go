package router

import (
	"github.com/gin-gonic/gin"

	"mxshop_api/order_web/api/order"
	"mxshop_api/order_web/middlewares"
)

func InitOrderRouter(Router *gin.RouterGroup) {
	orderRouter := Router.Group("orders").Use(middlewares.JWTAuth(),middlewares.Trace())
	{
		orderRouter.GET("",order.List)
		orderRouter.POST("",order.New)
		orderRouter.GET("/:id",order.Detail)
	}
}
