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
	"github.com/opentracing/opentracing-go"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

	"mxshop_srvs/order_srv/global"
	"mxshop_srvs/order_srv/handler"
	"mxshop_srvs/order_srv/initialize"
	"mxshop_srvs/order_srv/proto"
	"mxshop_srvs/order_srv/utils"
	"mxshop_srvs/order_srv/utils/otgrpc"
	"mxshop_srvs/order_srv/utils/register/consul"
)
func main() {
	IP := flag.String("ip","192.168.1.129","ip地址")
	Port := flag.Int("port",4444,"端口号")
	flag.Parse()

	if *Port == 50052 {
		*Port,_= utils.GetFreePort()
	}
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDb()
	initialize.InitRedsync()
	initialize.InitSrvs()
	initialize.InitProducer()

	serviceId := fmt.Sprintf("%s",uuid.NewV4())
	zap.S().Infof("uuid:%s,address:%s:%d",serviceId,*IP,*Port)


	//初始化tracer
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
			LocalAgentHostPort:fmt.Sprintf("%s:%d",global.ServerConfig.JaegerInfo.Host, global.ServerConfig.JaegerInfo.Port),
		},
		ServiceName: global.ServerConfig.JaegerInfo.Name,
	}

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		zap.S().Panic(err)
	}
	opentracing.SetGlobalTracer(tracer)

	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
	proto.RegisterOrderServer(server,&handler.OrderServer{})
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
		consumer.WithGroupName("mxshop_order"),
	)
	err = c.Subscribe("order_timeout", consumer.MessageSelector{}, handler.OrderTimeout)
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

	global.Producer.Shutdown()
	global.TransactionProduce.Shutdown()
	closer.Close()

	if err = registerClient.DeRegister(serviceId); err != nil {
		zap.S().Info("注销失败:", err.Error())
	}else{
		zap.S().Info("注销成功:")
	}
}
