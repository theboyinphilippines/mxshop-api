package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/proto"
)

func InitSrvConn() {
	// 负载均衡连接consul中的user-srv服务
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			global.ServerConfig.ConsulInfo.Host,
			global.ServerConfig.ConsulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接【用户服务失败】")
	}

	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}
// 客户端没用负载均衡连接consul的代码
func InitSrvConn2() {
	// 从注册中心获取到用户服务的信息
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d",
		global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)

	userSrvHost := ""
	userSrvPort := 0

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service==%s`, global.ServerConfig.UserSrvInfo.Name))
	if err != nil {
		return
	}
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		// 已经过滤好了，直接循环一次就可以break
		break
	}

	if userSrvHost == "" {
		zap.S().Fatal("[InitSrcConn2] 连接【用户服务失败】")
		return
	}

	//拨号连接用户grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接【用户服务】失败", "msg", err.Error())
	}

	//todo 后续服务下线，ip改变等，这里初始化会出现问题，留一个后续
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient

}
