package initialize

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/order-web/middlewares"
	"mxshop-api/order-web/router"
	"net/http"
)

func Routers() *gin.Engine {

	Router := gin.Default()
	//配置consul服务注册健康检查
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	// 配置跨域
	Router.Use(middlewares.Cors())
	ApiGroup := Router.Group("/o/v1")
	router.InitOrderRouter(ApiGroup)
	router.InitShopCartRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)
	return Router
}
