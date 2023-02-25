package router

import (
	"github.com/gin-gonic/gin"

	"mxshop_api/goods_web/api/category"
	"mxshop_api/goods_web/middlewares"
)

func InitCategoryRouter(Router *gin.RouterGroup) {
	CategoryRouter := Router.Group("categorys").Use(middlewares.Trace())
	{
		CategoryRouter.GET("", category.List)
		CategoryRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), category.New)
		CategoryRouter.GET("/:id", category.Detail)
		CategoryRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), category.Delete)
		CategoryRouter.PUT("/:id",middlewares.JWTAuth(), middlewares.IsAdminAuth(), category.Update)
	}
}
