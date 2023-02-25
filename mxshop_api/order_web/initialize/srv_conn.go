package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mxshop_api/order_web/global"
	"mxshop_api/order_web/proto"
	"mxshop_api/order_web/utils/otgrpc"
)

func InitSrvConn() {
	goodsConn,err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",global.ServerConfig.ConsulInfo.Host,global.ServerConfig.ConsulInfo.Port,global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【商品服务失败】",
			"msg", err.Error(),
		)
	}

	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)

	orderConn,err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",global.ServerConfig.ConsulInfo.Host,global.ServerConfig.ConsulInfo.Port,global.ServerConfig.OrderSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【订单服务失败】",
			"msg", err.Error(),
		)
	}

	global.OrderSrvClient = proto.NewOrderClient(orderConn)

	inventoryConn,err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",global.ServerConfig.ConsulInfo.Host,global.ServerConfig.ConsulInfo.Port,global.ServerConfig.InventorySrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【订单服务失败】",
			"msg", err.Error(),
		)
	}

	global.InventorySrvClient = proto.NewInventoryClient(inventoryConn)
}

