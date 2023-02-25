package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/handler"
	"mxshop_srvs/goods_srv/initialize"
	"mxshop_srvs/goods_srv/proto"
	"mxshop_srvs/goods_srv/utils"
)
func main() {
	IP := flag.String("ip","192.168.1.129","ip地址")
	Port := flag.Int("port",6666,"端口号")
	flag.Parse()

	if *Port == 50052 {
		*Port,_= utils.GetFreePort()
	}
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDb()
	initialize.InitEs()

	uuid := fmt.Sprintf("%s",uuid.NewV4())
	zap.S().Infof("uuid:%s,address:%s:%d",uuid,*IP,*Port)

	server := grpc.NewServer()
	proto.RegisterGoodsServer(server,&handler.GoodsServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("fail to listen" + err.Error())
	}

	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server,health.NewServer())
	//服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	client,err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name:    global.ServerConfig.Name,
		ID:      uuid,
		Port:    *Port,
		Tags:    global.ServerConfig.Tags,
		Address: global.ServerConfig.Host,
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port),
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
	if err != nil {
		panic("fail to register service" + err.Error())
	}

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("fail to start grpc" + err.Error())
		}
	}()

	//接收终止信号立即进行服务注销
	quit := make(chan os.Signal)
	signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	<-quit
	err = client.Agent().ServiceDeregister(uuid)
	if err != nil {
		zap.S().Info("注销失败")
	}else{
		zap.S().Info("注销成功:")
	}
}
