package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
	"mxshop-api/user-web/utils"
	"mxshop-api/user-web/utils/register/consul"
	myvalidator "mxshop-api/user-web/validator"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	//logger, _ := zap.NewProduction() //生产环境
	////logger, _ := zap.NewDevelopment() //开发环境
	//defer logger.Sync() // flushes buffer, if any
	//url := "https://imooc.com"
	//sugar := logger.Sugar()
	//sugar.Infow("failed to fetch URL",
	//	// Structured context as loosely typed key-value pairs.
	//	"url", url,
	//	"attempt", 3,
	//	"backoff", time.Second,
	//)
	//sugar.Infof("Failed to fetch URL: %s", url)

	// 1.初始化logger
	initialize.InitLogger()

	// 2. 初始化配置文件
	initialize.InitConfig()

	//3. 初始化routers
	Router := initialize.Routers()

	//4.初始化翻译器
	_ = initialize.InitTrans("zh")

	//5.初始化连接consul服务
	initialize.InitSrvConn()

	// 注册自定义验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			//基于自定义mobile标签 创建一个翻译器
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	//用zap.S()代替

	//logger, _ := zap.NewProduction()
	//defer logger.Sync()
	//super := logger.Sugar()

	// S()可以获取一个全局的sugar, 可以让我们自己设置一个全局的logger
	// S()函数和L()函数很有用: 提供了一个全局的安全访问logger的途径

	// 生产环境用动态可用端口
	viper.AutomaticEnv()
	debug := viper.GetBool("MZSHOP_DEBUG")
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}
	//服务注册到consul中
	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err := registerClient.Register(global.ServerConfig.Host,
		global.ServerConfig.Port,
		global.ServerConfig.Name,
		global.ServerConfig.Tags,
		serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败：", err.Error())
	}

	zap.S().Debugf("启动服务器，端口：%d", global.ServerConfig.Port)

	go func() {
		err = Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port))
		if err != nil {
			zap.S().Panic("启动失败：", err.Error())
		}
	}()

	//优雅退出，注销consul中服务
	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	err = registerClient.DeRegister(serviceId)
	if err != nil {
		zap.S().Info("注销失败：", err.Error())
	} else {
		zap.S().Info("注销成功")
	}
}
