package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mxshop_api/goods_web/global"
	"mxshop_api/goods_web/proto"
	"mxshop_api/goods_web/utils/otgrpc"
)

func InitSrvConn() {
	userConn,err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",global.ServerConfig.ConsulInfo.Host,global.ServerConfig.ConsulInfo.Port,global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】",
			"msg", err.Error(),
		)
	}

	global.GoodsSrvClient = proto.NewGoodsClient(userConn)
}

