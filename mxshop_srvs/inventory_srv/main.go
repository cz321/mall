package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"mxshop_srvs/inventory_srv/global"
	"mxshop_srvs/inventory_srv/handler"
	"mxshop_srvs/inventory_srv/initialize"
	"mxshop_srvs/inventory_srv/proto"
	"mxshop_srvs/inventory_srv/utils"
	"mxshop_srvs/inventory_srv/utils/register/consul"
)
func main() {
	IP := flag.String("ip","192.168.1.129","ip地址")
	Port := flag.Int("port",5555,"端口号")
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
	proto.RegisterInventoryServer(server,&handler.InventoryServer{})
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

	//监听库存归还
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"192.168.10.105:9876"}),
		consumer.WithGroupName("mxshop_inventory"),
	)
	err = c.Subscribe("order_reback", consumer.MessageSelector{}, handler.AutoReback)
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	_ = c.Start()


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

	_ = c.Shutdown()

	if err = registerClient.DeRegister(serviceId); err != nil {
		zap.S().Info("注销失败:", err.Error())
	}else{
		zap.S().Info("注销成功:")
	}
}
