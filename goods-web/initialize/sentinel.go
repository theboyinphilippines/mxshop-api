package initialize

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
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

	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		// Statistic time span=5s, recoveryTimeout=3s, maxErrorCount=50
		{
			Resource:         "abc",
			Strategy:         circuitbreaker.ErrorCount,
			RetryTimeoutMs:   3000, //3s之后尝试恢复
			MinRequestAmount: 10,   //静默数，10个请求以内全部通过
			StatIntervalMs:   5000, //5s统计一次
			Threshold:        50,   //错误数不超过50个
		},
	})
	if err != nil {
		log.Fatal(err)
	}

}
