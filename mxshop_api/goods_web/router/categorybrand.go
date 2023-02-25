package router

import (
	"github.com/gin-gonic/gin"

	"mxshop_api/goods_web/api/categoryBrand"
	"mxshop_api/goods_web/middlewares"
)

func InitCategoryBrandRouter(Router *gin.RouterGroup) {
	CategoryBrandRouter := Router.Group("categorybrands").Use(middlewares.Trace())
	{
		CategoryBrandRouter.GET("", categoryBrand.List)	//获取所有品牌
		CategoryBrandRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), categoryBrand.New)
		CategoryBrandRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), categoryBrand.Update)
		CategoryBrandRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), categoryBrand.Delete)
		CategoryBrandRouter.GET("/:id", categoryBrand.GetCategoryBrandList) //获取指定分类的所有品牌
	}
}
