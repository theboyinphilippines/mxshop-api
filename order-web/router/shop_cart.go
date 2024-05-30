package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/order-web/api/shop_cart"
	"mxshop-api/order-web/middlewares"
)

func InitShopCartRouter(Router *gin.RouterGroup) {
	ShopCartRouter := Router.Group("/shopcarts").Use(middlewares.JWTAuth())
	zap.S().Info("配置购物车相关的日志")
	{
		ShopCartRouter.GET("", shop_cart.List)          //购物车列表
		ShopCartRouter.POST("", shop_cart.New)          //添加商品到购物车
		ShopCartRouter.PATCH("/:id", shop_cart.Update)  // 修改条目
		ShopCartRouter.DELETE("/:id", shop_cart.Delete) //删除条目
	}

}
