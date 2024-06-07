package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mxshop-api/userop-web/global"
	"mxshop-api/userop-web/proto"
)

// 从consul中服务发现
func InitSrvConn() {
	// 负载均衡连接consul中的userop-srv服务
	useropConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			global.ServerConfig.ConsulInfo.Host,
			global.ServerConfig.ConsulInfo.Port, global.ServerConfig.UseropSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接【用户操作服务失败】")
	}

	global.UserFavSrvClient = proto.NewUserFavClient(useropConn)
	global.MessageSrvClient = proto.NewMessageClient(useropConn)
	global.AddressSrvClient = proto.NewAddressClient(useropConn)

	// 负载均衡连接consul中的goods-srv服务
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			global.ServerConfig.ConsulInfo.Host,
			global.ServerConfig.ConsulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接【商品服务失败】")
	}

	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)
}
