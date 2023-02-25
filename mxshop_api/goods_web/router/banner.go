package router

import (
	"github.com/gin-gonic/gin"

	"mxshop_api/goods_web/api/banner"
	"mxshop_api/goods_web/middlewares"
)

func InitBannerRouter(Router *gin.RouterGroup) {
	BannerRouter := Router.Group("banners").Use(middlewares.Trace())
	{
		BannerRouter.GET("",banner.List)
		BannerRouter.POST("",middlewares.JWTAuth(),middlewares.IsAdminAuth(),banner.New)
		BannerRouter.PUT("/:id",middlewares.JWTAuth(),middlewares.IsAdminAuth(),banner.Update)
		BannerRouter.DELETE("/:id",middlewares.JWTAuth(),middlewares.IsAdminAuth(),banner.Delete)
	}
}
