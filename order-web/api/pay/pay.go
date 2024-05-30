package pay

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 第三方库 需要go 1.8
func Notify(ctx *gin.Context) {
	//支付宝回调通知（处理notify_url带来的信息）
	//client, err := alipay.New(global.ServerConfig.AliPayInfo.AppID, global.ServerConfig.AliPayInfo.PrivateKey, false)
	//if err != nil {
	//	zap.S().Errorw("实例化支付宝失败")
	//	ctx.JSON(http.StatusInternalServerError, gin.H{
	//		"msg": err.Error(),
	//	})
	//	return
	//}
	//err = client.LoadAliPayPublicKey((global.ServerConfig.AliPayInfo.AliPublicKey))
	//if err != nil {
	//	zap.S().Errorw("加载支付宝的公钥失败")
	//	ctx.JSON(http.StatusInternalServerError, gin.H{
	//		"msg": err.Error(),
	//	})
	//	return
	//}
	//// DecodeNotification 内部已调用 VerifySign 方法验证签名
	//noti, err := client.DecodeNotification(ctx.Request.Form)
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, gin.H{
	//		"msg": err.Error(),
	//	})
	//	return
	//}
	//// 业务处理（修改订单号状态：拿到支付宝传过来的 商户订单号和订单状态， 写入到订单表）
	//_, err = global.OrderSrvClient.UpdateOrderStatus(context.Background(), &proto.OrderStatus{
	//	OrderSn: noti.OutTradeNo,
	//	Status:  string(noti.TradeStatus),
	//})
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, gin.H{})
	//	return
	//}
	// 如果通知消息没有问题，我们需要确认收到通知消息，不然支付宝后续会继续推送相同的消息
	//alipay.ACKNotification(ctx.Writer)
	ctx.String(http.StatusOK, "success")
}
