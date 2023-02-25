package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mxshop_api/user_web/middlewares"
	"mxshop_api/user_web/router"
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

	apiGroup := engine.Group("/u/v1")
	router.InitUserRouter(apiGroup)
	router.InitBaseRouter(apiGroup)

	return engine
}