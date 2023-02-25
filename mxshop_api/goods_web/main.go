package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"

	"mxshop_api/goods_web/global"
	"mxshop_api/goods_web/initialize"
	"mxshop_api/goods_web/utils"
	"mxshop_api/goods_web/utils/consul"
)

func main() {
	//初始化logger
	initialize.InitLogger()
	//初始化配置文件
	initialize.InitConfig()
	//初始化routers
	router := initialize.Routers()
	//初始化翻译器
	initialize.InitTrans("zh")
	//初始化srv_conn
	initialize.InitSrvConn()
	//初始化sentinel
	initialize.InitSentinel()

	debug := initialize.GetEnvInfo("MXSHOP_DEBUG")
	if !debug {
		port,err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err := registerClient.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
	}
	zap.S().Debugf("启动服务器, 端口： %d", global.ServerConfig.Port)
	go func(){
		if err := router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil{
			zap.S().Panic("启动失败:", err.Error())
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
