package middlewares

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/user_web/models"
	"net/http"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		curUser := claims.(*models.CustomClaims)
		if curUser.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden,gin.H{
				"msg":"没有权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
