package initialize

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"

	"mxshop_api/goods_web/global"
)

func InitConfig() {
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("goods_web/%s_pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("goods_web/%s_debug.yaml", configFilePrefix)
	}

	v := viper.New()

	v.SetConfigFile(configFileName)

	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = v.Unmarshal(global.NacosConfig)
	if err != nil {
		panic(err)
	}
	zap.S().Infof("nacos配置信息：%v", *global.NacosConfig)

	getConfigFromNacos()

	zap.S().Infof("server配置信息：%v", *global.ServerConfig)
}

//从Nacos中读取配置信息
func getConfigFromNacos() {
	// 创建ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}

	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// 创建动态配置客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		zap.S().Fatalf("创建动态配置客户端错误:%s", err.Error())
	}
	// 获取配置
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group})

	if err != nil {
		zap.S().Fatalf("获取配置错误:%s", err.Error())
	}

	err = json.Unmarshal([]byte(content), global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("json.Unmarshal错误:%s", err.Error())
	}
}

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}
