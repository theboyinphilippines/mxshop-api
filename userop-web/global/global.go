package global

import (
	ut "github.com/go-playground/universal-translator"
	"mxshop-api/userop-web/config"
	"mxshop-api/userop-web/proto"
)

var (
	Trans            ut.Translator
	ServerConfig     *config.ServerConfig = &config.ServerConfig{}
	GoodsSrvClient   proto.GoodsClient
	MessageSrvClient proto.MessageClient
	AddressSrvClient proto.AddressClient
	UserFavSrvClient proto.UserFavClient
	NacosConfig      *config.NacosConfig = &config.NacosConfig{}
)
