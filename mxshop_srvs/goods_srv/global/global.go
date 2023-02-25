package global

import (
	"gorm.io/gorm"

	"github.com/olivere/elastic/v7"

	"mxshop_srvs/goods_srv/config"
)
var (
	DB *gorm.DB

	ServerConfig config.ServerConfig

	NacosConfig config.NacosConfig

	EsClient *elastic.Client
)