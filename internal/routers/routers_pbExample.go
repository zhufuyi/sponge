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

// nolint
var (
	apiV1RouterFns_pbExample []func(prePath string, engine *gin.RouterGroup) // group router functions
	// if you have other group routes you can define them here
	// example:
	//     myPrePathRouterFns []func(prePath string, engine *gin.RouterGroup)
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
	if config.Get().App.EnableTrace {
		r.Use(middleware.Tracing(config.Get().App.Name))
	}

	// pprof performance analysis
	if config.Get().App.EnableHTTPProfile {
		prof.Register(r, prof.WithIOWaitTime())
	}

	// validator
	binding.Validator = validator.Init()

	r.GET("/health", handlerfunc.CheckHealth)
	r.GET("/ping", handlerfunc.Ping)

	// access path /apis/swagger/index.html
	swagger.CustomRouter(r, "apis", docs.ApiDocs)

	// register routers, middleware support
	registerRouters_pbExample(r, "/api/v1", apiV1RouterFns_pbExample)
	// if you have other group routes you can add them here
	// example:
	//    registerRouters(r, "/myPrePath", myPrePathRouterFns, middleware.Auth())

	return r
}

// nolint
func registerRouters_pbExample(r *gin.Engine, prePath string,
	routerFns []func(prePath string, engine *gin.RouterGroup), handlers ...gin.HandlerFunc) {
	rg := r.Group(prePath, handlers...)
	for _, fn := range routerFns {
		fn(prePath, rg)
	}
}
