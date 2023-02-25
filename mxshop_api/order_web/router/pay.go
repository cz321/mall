package router

import (
	"github.com/gin-gonic/gin"

	"mxshop_api/order_web/api/pay"
)

func InitPayRouter(Router *gin.RouterGroup) {
	orderRouter := Router.Group("pay")
	{
		orderRouter.POST("alipay/notify", pay.Notify)
	}
}
