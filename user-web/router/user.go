package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/user-web/api"
	"mxshop-api/user-web/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("/user")
	zap.S().Info("配置用户相关的日志")
	{
		UserRouter.GET("/list",middlewares.JWTAuth(),middlewares.IsAdminAuth(),api.GetUserList)
		UserRouter.POST("/pwd_login",api.PasswordLogin)
		UserRouter.POST("/register",api.Register)
	}

}