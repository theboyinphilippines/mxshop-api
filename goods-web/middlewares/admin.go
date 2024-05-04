package middlewares

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/goods-web/models"
	"net/http"
)

// 判断是否是管理员权限的接口 中间件 管理员角色id为2
func IsAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims,_ := c.Get("claims")
		currentUser := claims.(*models.CustomClaims)
		if currentUser.AuthorityId != 2 {
			c.JSON(http.StatusForbidden,gin.H{
				"msg": "无权限",
			})
			c.Abort()
			return
		}
		c.Next()

	}

}
