package initialize

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mxshop_api/user_web/global"
	"mxshop_api/user_web/proto"
)

func InitSrvConn() {
	userConn,err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",global.ServerConfig.ConsulInfo.Host,global.ServerConfig.ConsulInfo.Port,global.ServerConfig.UserSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】",
			"msg", err.Error(),
		)
	}

	global.UserSrvClient = proto.NewUserClient(userConn)
}
func InitSrvConn1()  {
	//服务发现
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d",global.ServerConfig.ConsulInfo.Host,global.ServerConfig.ConsulInfo.Port)

	client,err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//服务过滤
	services, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`,global.ServerConfig.UserSrvInfo.Name))

	if err != nil {
		panic(err)
	}

	userSrvHost := ""
	userSrvPort := 0
	for _,value := range services {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}

	if userSrvHost == "" {
		zap.S().Fatal("【InitSrcConn】 获取用户服务失败")
	}
	zap.S().Infof("userSrvAddress: %s:%d",userSrvHost,userSrvPort)

	//拨号
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】",
			"msg", err.Error(),
		)
	}

	global.UserSrvClient = proto.NewUserClient(userConn)
}

