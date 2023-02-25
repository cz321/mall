package router

import (
	"github.com/gin-gonic/gin"

	"mxshop_api/order_web/api/shop_cart"
	"mxshop_api/order_web/middlewares"
)

func InitShopCartRouter(Router *gin.RouterGroup) {
	shopCartRouter := Router.Group("shopcarts").Use(middlewares.JWTAuth())
	{
		shopCartRouter.GET("", shop_cart.List)
		shopCartRouter.DELETE("/:id", shop_cart.Delete)
		shopCartRouter.POST("", shop_cart.New)
		shopCartRouter.PATCH("/:id",shop_cart.Update)
	}
}
