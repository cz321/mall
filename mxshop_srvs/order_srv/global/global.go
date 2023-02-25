package global

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"

	"mxshop_srvs/order_srv/config"
	"mxshop_srvs/order_srv/proto"
)

var (
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig

	DB      *gorm.DB
	Redsync *redsync.Redsync

	GoodsSrvClient     proto.GoodsClient
	InventorySrvClient proto.InventoryClient

	TransactionProduce rocketmq.TransactionProducer
	Producer           rocketmq.Producer
)
