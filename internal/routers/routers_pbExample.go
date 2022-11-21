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
	rootRouterFns []func(engine *gin.Engine) // root routing group, used by rpc gateway
)

// NewRouter_pbExample create a new router
func NewRouter_pbExample() *gin.Engine { //nolint
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	// request id middleware
	r.Use(middleware.RequestID())

	// logger middleware
	r.Use(middleware.Logging(
		middleware.WithLog(logger.Get()),
		middleware.WithRequestIDFromContext(),
		middleware.WithIgnoreRoutes("/metrics"), // ignore path
	))

	// metrics middleware
	if config.Get().App.EnableMetrics {
		r.Use(metrics.Metrics(r,
			//metrics.WithMetricsPath("/metrics"),                // default is /metrics
			metrics.WithIgnoreStatusCodes(http.StatusNotFound), // ignore 404 status codes
		))
	}

	// limit middleware
	if config.Get().App.EnableLimit {
		r.Use(middleware.RateLimit())
	}

	// circuit breaker middleware
	if config.Get().App.EnableCircuitBreaker {
		r.Use(middleware.CircuitBreaker())
	}

	// trace middleware
	if config.Get().App.EnableTracing {
		r.Use(middleware.Tracing(config.Get().App.Name))
	}

	// pprof performance analysis
	if config.Get().App.EnablePprof {
		prof.Register(r, prof.WithIOWaitTime())
	}

	// validator
	binding.Validator = validator.Init()

	r.GET("/health", handlerfunc.CheckHealth)
	r.GET("/ping", handlerfunc.Ping)

	// access path /apis/swagger/index.html
	swagger.CustomRouter(r, "apis", docs.ApiDocs)

	// registration/Prefix Routing Groups
	for _, fn := range rootRouterFns {
		fn(r)
	}

	return r
}
