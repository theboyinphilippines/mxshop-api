package global

import (
	ut "github.com/go-playground/universal-translator"
	"mxshop-api/order-web/config"
	"mxshop-api/order-web/proto"
)

var (
	Trans              ut.Translator
	ServerConfig       *config.ServerConfig = &config.ServerConfig{}
	OrderSrvClient     proto.OrderClient
	GoodsSrvClient     proto.GoodsClient
	InventorySrvClient proto.InventoryClient
	NacosConfig        *config.NacosConfig = &config.NacosConfig{}
)
