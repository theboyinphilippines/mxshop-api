package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := jaegercfg.Configuration{
			//采样器设置
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			//jaeger agent设置
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans:           true,
				LocalAgentHostPort: "192.168.0.101:6831",
			},
			ServiceName: "mxshop-goods-api",
		}
		tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
		if err != nil {
			panic(any(err))
		}
		opentracing.SetGlobalTracer(tracer)
		defer closer.Close()
		parentSpan := tracer.StartSpan(c.Request.URL.Path)
		defer parentSpan.Finish()
		c.Set("tracer", tracer)
		c.Set("parentSpan", parentSpan)
		c.Next()
	}

}
