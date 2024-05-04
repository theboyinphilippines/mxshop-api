package initialize

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/goods-web/middlewares"
	router2 "mxshop-api/goods-web/router"
)

func Routers() *gin.Engine {

	Router := gin.Default()
	// 配置跨域
	Router.Use(middlewares.Cors())
	ApiGroup := Router.Group("/g/v1")
	router2.InitGoodsRouter(ApiGroup)
	router2.InitBaseRouter(ApiGroup)
	return Router
}
