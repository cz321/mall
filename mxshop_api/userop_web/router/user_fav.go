package router

import (
	"github.com/gin-gonic/gin"

	"mxshop_api/userop_web/api/user_fav"
	"mxshop_api/userop_web/middlewares"
)

func InitUserFavRouter(Router *gin.RouterGroup) {
	userFavRouter := Router.Group("userfavs").Use(middlewares.JWTAuth())
	{
		userFavRouter.GET("", user_fav.List)
		userFavRouter.DELETE("/:id", user_fav.Delete)
		userFavRouter.POST("", user_fav.New)
		userFavRouter.GET("/:id",user_fav.Detail)
	}
}
