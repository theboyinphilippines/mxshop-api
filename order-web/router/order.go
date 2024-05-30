package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/order-web/api/order"
	"mxshop-api/order-web/api/pay"
)

func InitOrderRouter(Router *gin.RouterGroup) {
	OrderRouter := Router.Group("/orders")
	zap.S().Info("配置订单相关的日志")
	{
		OrderRouter.GET("", order.List)        //购物车列表
		OrderRouter.POST("", order.New)        //新建订单
		OrderRouter.POST("/:id", order.Detail) //订单详情
	}
	PayRouter := Router.Group("pay")
	{
		PayRouter.POST("alipay/notify", pay.Notify)
	}

}
