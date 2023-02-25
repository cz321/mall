package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mxshop_api/userop_web/global"
	"mxshop_api/userop_web/proto"
)

func InitSrvConn() {
	goodsConn,err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",global.ServerConfig.ConsulInfo.Host,global.ServerConfig.ConsulInfo.Port,global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【商品服务失败】",
			"msg", err.Error(),
		)
	}

	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)

	userOpConn,err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",global.ServerConfig.ConsulInfo.Host,global.ServerConfig.ConsulInfo.Port,global.ServerConfig.UserOpSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户操作服务失败】",
			"msg", err.Error(),
		)
	}

	global.AddressClient = proto.NewAddressClient(userOpConn)
	global.MessageSrvClient = proto.NewMessageClient(userOpConn)
	global.UserFavSrvClient = proto.NewUserFavClient(userOpConn)
}

