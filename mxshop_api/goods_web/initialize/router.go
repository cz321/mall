package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mxshop_api/goods_web/middlewares"
	"mxshop_api/goods_web/router"
)

func Routers() *gin.Engine {
	engine := gin.Default()

	//配置跨域请求
	engine.Use(middlewares.Cors())

	//健康检查
	engine.GET("/health", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"code":http.StatusOK,
			"success":true,
		})
	})

	apiGroup := engine.Group("/g/v1")

	router.InitGoodsRouter(apiGroup)	//商品
	router.InitCategoryRouter(apiGroup)	//商品分类
	router.InitBannerRouter(apiGroup)	//轮播图
	router.InitBrandRouter(apiGroup)	//品牌
	router.InitCategoryBrandRouter(apiGroup)	//品牌分类

	return engine
}