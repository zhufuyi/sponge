package routers

import (
	"net/http"
	"strings"

	"github.com/zhufuyi/sponge/docs"
	"github.com/zhufuyi/sponge/internal/config"
	"github.com/zhufuyi/sponge/pkg/gin/handlerfunc"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/middleware/metrics"
	"github.com/zhufuyi/sponge/pkg/gin/middleware/ratelimiter"
	"github.com/zhufuyi/sponge/pkg/gin/validator"
	"github.com/zhufuyi/sponge/pkg/logger"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	routerFns []func()         // 路由集合
	apiV1     *gin.RouterGroup // 基础路由组
)

// NewRouter 实例化路由
func NewRouter() *gin.Engine {
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
		opts := []ratelimiter.Option{
			ratelimiter.WithQPS(config.Get().RateLimiter.QPSLimit),
			ratelimiter.WithBurst(config.Get().RateLimiter.MaxLimit),
		}
		if strings.ToUpper(config.Get().RateLimiter.Dimension) == "IP" {
			opts = append(opts, ratelimiter.WithIP())
		}
		r.Use(ratelimiter.QPS(opts...))
	}

	// trace 中间件
	if config.Get().App.EnableTracing {
		r.Use(middleware.Tracing(config.Get().App.Name))
	}

	// profile 性能分析
	if config.Get().App.EnableProfile {
		pprof.Register(r)
	}

	// 校验器
	binding.Validator = validator.Init()

	// 注册swagger路由，通过swag init生成代码
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", handlerfunc.CheckHealth)
	r.GET("/ping", handlerfunc.Ping)

	apiV1 = r.Group("/api/v1")

	// 注册所有路由
	for _, fn := range routerFns {
		fn()
	}

	return r
}
