package initialize

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"

	"mxshop_srvs/order_srv/global"
	"mxshop_srvs/order_srv/handler"
)

func InitProducer() {
	orderListener := handler.NewOrderListener(nil)

	global.TransactionProduce, _ = rocketmq.NewTransactionProducer(
		orderListener,
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"192.168.10.105:9876"})),
		producer.WithInstanceName("order_create"),
		producer.WithRetry(1),
	)

	//启动producer
	err := global.TransactionProduce.Start()
	if err != nil {
		zap.S().Panic("启动TransactionProduce失败", err)
	}


	//发送延时消息
	global.Producer, _ = rocketmq.NewProducer(
		producer.WithNameServer([]string{"192.168.10.105:9876"}),
		producer.WithInstanceName("order_others"),

	)

	err = global.Producer.Start()
	if err != nil {
		zap.S().Panic("生成Product失败",err.Error())
	}
}
