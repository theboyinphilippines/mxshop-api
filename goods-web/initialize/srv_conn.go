package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/proto"
)

// 从consul中服务发现
func InitSrvConn() {
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

	goodsSrvClient := proto.NewGoodsClient(goodsConn)
	global.GoodsSrvClient = goodsSrvClient
}

// 客户端没用负载均衡连接consul的代码
func InitSrvConn2() {
	// 从注册中心获取到用户服务的信息
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d",
		global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)

	goodsSrvHost := ""
	goodsSrvPort := 0

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service==%s`, global.ServerConfig.GoodsSrvInfo.Name))
	if err != nil {
		return
	}
	for _, value := range data {
		goodsSrvHost = value.Address
		goodsSrvPort = value.Port
		// 已经过滤好了，直接循环一次就可以break
		break
	}

	if goodsSrvHost == "" {
		zap.S().Fatal("[InitSrcConn2] 连接【用户服务失败】")
		return
	}

	//拨号连接用户grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", goodsSrvHost, goodsSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetGoodsList] 连接【用户服务】失败", "msg", err.Error())
	}

	//todo 后续服务下线，ip改变等，这里初始化会出现问题，留一个后续
	goodsSrvClient := proto.NewGoodsClient(userConn)
	global.GoodsSrvClient = goodsSrvClient

}
