package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/userop-web/api/message"
	"mxshop-api/userop-web/middlewares"
)

func InitMessageRouter(Router *gin.RouterGroup) {
	MessageRouter := Router.Group("message").Use(middlewares.JWTAuth())
	{
		MessageRouter.GET("", message.List) // 留言列表页
		MessageRouter.POST("", message.New) //新建留言
	}
}
