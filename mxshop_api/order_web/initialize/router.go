package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mxshop_api/order_web/middlewares"
	"mxshop_api/order_web/router"
)

func Routers() *gin.Engine {
	engine := gin.Default()

	//配置跨域请求
	engine.Use(middlewares.Cors())

	//健康检查
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	apiGroup := engine.Group("/o/v1")

	router.InitOrderRouter(apiGroup)    //订单
	router.InitShopCartRouter(apiGroup) //购物车
	router.InitPayRouter(apiGroup)      //支付

	return engine
}
