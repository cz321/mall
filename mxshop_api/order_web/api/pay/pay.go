package pay

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"

	"mxshop_api/order_web/api"
	"mxshop_api/order_web/global"
	"mxshop_api/order_web/proto"
)


func Notify(ctx *gin.Context) {
	alipayInfo := global.ServerConfig.AlipayInfo
	var client, err = alipay.New(alipayInfo.AppID,alipayInfo.PrivateKey, false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"msg":"实例化 AliPay 失败",
		})
		return
	}
	err = client.LoadAliPayPublicKey(alipayInfo.AliPublicKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"msg":"加载AliPulickey失败",
		})
		return
	}
	noti, err := client.GetTradeNotification(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"msg":"不合法的通知",
		})
		return
	}
	_, err = global.OrderSrvClient.UpdateOrderStatus(context.Background(), &proto.OrderStatus{
		Id:      0,
		OrderSn: noti.OutTradeNo,
		Status:  string(noti.TradeStatus),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err,ctx)
	}

	// 确认收到通知消息
	zap.S().Info("支付完毕")
	ctx.String(http.StatusOK,"success")
}


func GeneratePayUrl(OrderSn string, totalAmount float32) (string,error){
	return generateAliPayUrl(OrderSn,totalAmount)
}


func generateAliPayUrl (OrderSn string,totalAmount float32) (string,error) {
	alipayInfo := global.ServerConfig.AlipayInfo
	var client, err = alipay.New(alipayInfo.AppID,alipayInfo.PrivateKey, false)
	if err != nil {
		return "",nil
	}
	client.LoadAliPayPublicKey(alipayInfo.AliPublicKey)

	var p = alipay.TradePagePay{}
	p.NotifyURL = alipayInfo.NotifyURL
	p.ReturnURL = alipayInfo.ReturnURL
	p.Subject = "mxshop订单-" + OrderSn
	p.OutTradeNo = OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(totalAmount),'f',2,64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		return "",nil
	}
	return url.String(),nil
}

func generateWechatPayUrl (OrderSn string,totalAmount float32) (string,error) {
	return "",nil
}