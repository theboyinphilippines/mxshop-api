package global

import (
	ut "github.com/go-playground/universal-translator"
	"mxshop-api/goods-web/config"
	"mxshop-api/goods-web/proto"
)

var (
	Trans        ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	GoodsSrvClient proto.GoodsClient
	NacosConfig *config.NacosConfig  = &config.NacosConfig{}
)
