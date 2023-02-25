package router

import (
	"github.com/gin-gonic/gin"

	"mxshop_api/goods_web/api/goods"
	"mxshop_api/goods_web/middlewares"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods").Use(middlewares.Trace())
	{
		GoodsRouter.GET("", goods.List)
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New) //商品列表
		GoodsRouter.GET("/:id", goods.Detail)                                             //获取商品详情
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete)
		GoodsRouter.GET("/:id/stocks", goods.Stocks)
		GoodsRouter.PATCH("/:id",middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus)
		GoodsRouter.PUT("/:id",middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)
	}
}
