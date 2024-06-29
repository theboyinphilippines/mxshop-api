package initialize

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"log"
)

func InitSentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "goods-list",
			TokenCalculateStrategy: flow.Direct, //qps
			ControlBehavior:        flow.Reject, //直接拒绝
			Threshold:              3,
			StatIntervalInMs:       8000, //8秒内最大请求数量为3个，超过3个限流
		},
		{
			Resource:               "goods-detail",
			TokenCalculateStrategy: flow.Direct, //qps
			ControlBehavior:        flow.Reject, //直接拒绝
			Threshold:              3,
			StatIntervalInMs:       6000, //6秒内最大请求数量为3个，超过3个限流
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}

	//for i := 0; i < 10; i++ {
	//	e, b := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
	//	if b != nil {
	//		fmt.Println("通过")
	//	} else {
	//		fmt.Println("限流了")
	//		e.Exit()
	//	}
	//}
}
