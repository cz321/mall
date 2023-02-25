package global

import (
	ut "github.com/go-playground/universal-translator"

	"mxshop_api/goods_web/config"
	"mxshop_api/goods_web/proto"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig *config.NacosConfig = &config.NacosConfig{}

	Trans ut.Translator

	GoodsSrvClient proto.GoodsClient
)