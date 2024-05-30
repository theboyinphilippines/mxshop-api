package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/goods-web/api/goods"
	"mxshop-api/goods-web/middlewares"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("/order")
	zap.S().Info("配置商品相关的日志")
	{
		GoodsRouter.GET("/list", goods.List)
		GoodsRouter.POST("/new", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)
		GoodsRouter.GET("/detail/:id", goods.Detail)
		GoodsRouter.DELETE("delete/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete)
		GoodsRouter.GET("/stocks/:id", goods.Stocks)
		GoodsRouter.PUT("update/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)
		GoodsRouter.PATCH("patch/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus)
	}

}
