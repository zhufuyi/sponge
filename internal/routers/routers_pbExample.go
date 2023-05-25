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

type routeFns = []func(r *gin.Engine, groupPathMiddlewares map[string][]gin.HandlerFunc, singlePathMiddlewares map[string][]gin.HandlerFunc)

// all route functions
var allRouteFns = make(routeFns, 0)

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

	c := newMiddlewareConfig()

	// set up the middleware for the route group, route group is left prefix rules,
	// it can be viewed in the register() function in the api/serverNameExample/v1/xxx_route.go file
	// example:
	//		c.setMiddlewaresForGroupPath("/api/v1", middleware.Auth())
	//		c.setMiddlewaresForGroupPath("/api/v2", middleware.Auth(), middleware.RateLimit())

	// set up single route middleware, route must be full path,
	// it can be viewed in the register() function in the api/serverNameExample/v1/xxx_route.go file
	// example:
	//		c.setMiddlewaresForSinglePath("/api/v1/userExample", middleware.Auth())
	//		c.setMiddlewaresForSinglePath("/api/v1/userExample/list", middleware.Auth(), middleware.RateLimit())

	// register all routes
	registerAllRoutes(r, c, allRouteFns)

	return r
}

func registerAllRoutes(r *gin.Engine, c *middlewareConfig, routeFns routeFns) {
	if c == nil {
		c = newMiddlewareConfig()
	}
	for _, fn := range routeFns {
		fn(r, c.groupPathMiddlewares, c.singlePathMiddlewares)
	}
}

type middlewareConfig struct {
	groupPathMiddlewares  map[string][]gin.HandlerFunc // middleware function corresponding to route group
	singlePathMiddlewares map[string][]gin.HandlerFunc // middleware functions corresponding to a single route
}

func newMiddlewareConfig() *middlewareConfig {
	return &middlewareConfig{
		groupPathMiddlewares:  make(map[string][]gin.HandlerFunc),
		singlePathMiddlewares: make(map[string][]gin.HandlerFunc),
	}
}

func (c *middlewareConfig) setMiddlewaresForGroupPath(groupPath string, handlers ...gin.HandlerFunc) { // nolint
	if groupPath == "" {
		return
	}
	if groupPath[0] != '/' {
		groupPath = "/" + groupPath
	}

	handlerFns, ok := c.groupPathMiddlewares[groupPath]
	if !ok {
		c.groupPathMiddlewares[groupPath] = handlers
		return
	}

	c.groupPathMiddlewares[groupPath] = append(handlerFns, handlers...)
}

func (c *middlewareConfig) setMiddlewaresForSinglePath(singlePath string, handlers ...gin.HandlerFunc) { // nolint
	if singlePath == "" {
		return
	}
	handlerFns, ok := c.singlePathMiddlewares[singlePath]
	if !ok {
		c.singlePathMiddlewares[singlePath] = handlers
		return
	}

	c.singlePathMiddlewares[singlePath] = append(handlerFns, handlers...)
}
