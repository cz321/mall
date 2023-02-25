package router

import (
	"github.com/gin-gonic/gin"

	"mxshop_api/goods_web/api/brand"
	"mxshop_api/goods_web/middlewares"
)

func InitBrandRouter(Router *gin.RouterGroup) {
	BrandRouter := Router.Group("brands").Use(middlewares.Trace())
	{
		BrandRouter.GET("",brand.List)
		BrandRouter.POST("",middlewares.JWTAuth(),middlewares.IsAdminAuth(),brand.New)
		BrandRouter.PUT("/:id",middlewares.JWTAuth(),middlewares.IsAdminAuth(),brand.Update)
		BrandRouter.DELETE("/:id",middlewares.JWTAuth(),middlewares.IsAdminAuth(),brand.Delete)
	}
}