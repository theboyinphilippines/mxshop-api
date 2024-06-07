package initialize

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/userop-web/middlewares"
	"mxshop-api/userop-web/router"
	"net/http"
)

func Routers() *gin.Engine {

	Router := gin.Default()
	//为gin服务配置consul服务注册健康检查
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	// 配置跨域
	Router.Use(middlewares.Cors())
	ApiGroup := Router.Group("/up/v1")
	router.InitAddressRouter(ApiGroup)
	router.InitMessageRouter(ApiGroup)
	router.InitUserFavRouter(ApiGroup)
	return Router
}
