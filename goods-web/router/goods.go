package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/goods-web/api/goods"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("/goods")
	zap.S().Info("配置用户相关的日志")
	{
		GoodsRouter.GET("/list",goods.List)
	}

}