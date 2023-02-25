package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"mxshop_srvs/userop_srv/global"
	"mxshop_srvs/userop_srv/handler"
	"mxshop_srvs/userop_srv/initialize"
	"mxshop_srvs/userop_srv/proto"
	"mxshop_srvs/userop_srv/utils"
	"mxshop_srvs/userop_srv/utils/register/consul"
)
func main() {
	IP := flag.String("ip","192.168.1.129","ip地址")
	Port := flag.Int("port",7777,"端口号")
	flag.Parse()

	if *Port == 50052 {
		*Port,_= utils.GetFreePort()
	}
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDb()
	initialize.InitRedsync()

	serviceId := fmt.Sprintf("%s",uuid.NewV4())
	zap.S().Infof("uuid:%s,address:%s:%d",serviceId,*IP,*Port)

	server := grpc.NewServer()
	proto.RegisterAddressServer(server,&handler.UserOpServer{})
	proto.RegisterMessageServer(server,&handler.UserOpServer{})
	proto.RegisterUserFavServer(server,&handler.UserOpServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("fail to listen" + err.Error())
	}

	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server,health.NewServer())
	//服务注册
	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	err = registerClient.Register(global.ServerConfig.Host, *Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
	}

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("fail to start grpc" + err.Error())
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = registerClient.DeRegister(serviceId); err != nil {
		zap.S().Info("注销失败:", err.Error())
	}else{
		zap.S().Info("注销成功:")
	}
}
