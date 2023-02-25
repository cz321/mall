package global

import (
	ut "github.com/go-playground/universal-translator"

	"mxshop_api/userop_web/config"
	"mxshop_api/userop_web/proto"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig *config.NacosConfig = &config.NacosConfig{}

	Trans ut.Translator

	MessageSrvClient proto.MessageClient

	AddressClient proto.AddressClient

	UserFavSrvClient proto.UserFavClient

	GoodsSrvClient proto.GoodsClient
)