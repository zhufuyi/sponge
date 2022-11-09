package routers

import (
	"net/http"

	"github.com/zhufuyi/sponge/docs"
	"github.com/zhufuyi/sponge/internal/config"

	"github.com/zhufuyi/sponge/pkg/gin/handlerfunc"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/middleware/metrics"
	"github.com/zhufuyi/sponge/pkg/gin/prof"
	"github.com/zhufuyi/sponge/pkg/gin/swagger"
	"github.com/zhufuyi/sponge/pkg/gin/validator"
	"github.com/zhufuyi/sponge/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var (
	rootRouterFns []func(engine *gin.Engine) // 根路由组，rpc gateway使用
)

// NewRouter_pbExample 创建一个路由
func NewRouter_pbExample() *gin.Engine { //nolint
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	// request id 中间件
	r.Use(middleware.RequestID())

	// logger 中间件
	r.Use(middleware.Logging(
		middleware.WithLog(logger.Get()),
		middleware.WithRequestIDFromContext(),
		middleware.WithIgnoreRoutes("/metrics"), // 忽略路由
	))

	// metrics 中间件
	if config.Get().App.EnableMetrics {
		r.Use(metrics.Metrics(r,
			//metrics.WithMetricsPath("/metrics"),                // 默认是 /metrics
			metrics.WithIgnoreStatusCodes(http.StatusNotFound), // 忽略404状态码
		))
	}

	// limit 中间件
	if config.Get().App.EnableLimit {
		r.Use(middleware.RateLimit())
	}

	// circuit breaker 中间件
	if config.Get().App.EnableCircuitBreaker {
		r.Use(middleware.CircuitBreaker())
	}

	// trace 中间件
	if config.Get().App.EnableTracing {
		r.Use(middleware.Tracing(config.Get().App.Name))
	}

	// pprof 性能分析
	if config.Get().App.EnablePprof {
		prof.Register(r, prof.WithIOWaitTime())
	}

	// 校验器
	binding.Validator = validator.Init()

	r.GET("/health", handlerfunc.CheckHealth)
	r.GET("/ping", handlerfunc.Ping)

	swagger.CustomRouter(r, "apis", docs.ApiDocs)

	// 注册/前缀路由组
	for _, fn := range rootRouterFns {
		fn(r)
	}

	return r
}
