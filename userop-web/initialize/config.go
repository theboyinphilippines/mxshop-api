package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/userop-web/global"
)

// 配置环境变量，根据环境变量来决定用开发还是生产的配置文件
func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)

}

func InitConfig() {

	debug := GetEnvInfo("MZSHOP_DEBUG")
	fmt.Printf("debug是：%v\n", debug)
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("userop-web/%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("userop-web/%s-debug.yaml", configFilePrefix)
	}
	zap.S().Infof("配置文件路径为：%v", configFileName)

	v := viper.New()
	v.SetConfigFile(configFileName)
	zap.S().Infof("配置文件对象为：%v", v)
	if err := v.ReadInConfig(); err != nil {
		panic(any(err))
		//zap.S().Errorf("配置文件错误为：%v",err)
	}
	// serverConfig对象，其他文件中也要使用配置，所以声明为全局变量
	//serverConfig := config.ServerConfig{}
	if err := v.Unmarshal(global.NacosConfig); err != nil {
		panic(any(err))
	}
	zap.S().Infof("配置信息：%v", global.NacosConfig)
	fmt.Printf("服务名称是：%v", v.Get("name"))

	// 从nacos中读取配置
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.NamespaceId, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(any(err))
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group})

	if err != nil {
		panic(any(err))
	}

	//将从nacos中获取的配置数据绑定到结构体中
	fmt.Println("这是content", content)
	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("读取nacos配置失败： %s", err.Error())
	}
	fmt.Println("这是global.ServerConfig", &global.ServerConfig)

}
