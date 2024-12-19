package routers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/go-dev-frame/sponge/pkg/errcode"
	"github.com/go-dev-frame/sponge/pkg/gin/handlerfunc"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware/metrics"
	"github.com/go-dev-frame/sponge/pkg/gin/prof"
	"github.com/go-dev-frame/sponge/pkg/gin/swagger"
	"github.com/go-dev-frame/sponge/pkg/gin/validator"
	"github.com/go-dev-frame/sponge/pkg/jwt"
	"github.com/go-dev-frame/sponge/pkg/logger"

	"github.com/go-dev-frame/sponge/docs"
	"github.com/go-dev-frame/sponge/internal/config"
)

type routeFns = []func(r *gin.Engine, groupPathMiddlewares map[string][]gin.HandlerFunc, singlePathMiddlewares map[string][]gin.HandlerFunc)

var (
	// all route functions
	allRouteFns = make(routeFns, 0)
	// all middleware functions
	allMiddlewareFns = []func(c *middlewareConfig){}
)

// NewRouter_pbExample create a new router
func NewRouter_pbExample() *gin.Engine { //nolint
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	if config.Get().HTTP.Timeout > 0 {
		// if you need more fine-grained control over your routes, set the timeout in your routes, unsetting the timeout globally here.
		r.Use(middleware.Timeout(time.Second * time.Duration(config.Get().HTTP.Timeout)))
	}

	// request id middleware
	r.Use(middleware.RequestID())

	// logger middleware, to print simple messages, replace middleware.Logging with middleware.SimpleLog
	r.Use(middleware.Logging(
		middleware.WithLog(logger.Get()),
		middleware.WithRequestIDFromContext(),
		middleware.WithIgnoreRoutes("/metrics"), // ignore path
	))

	// init jwt middleware, you can replace it with your own jwt middleware
	jwt.Init(
	//jwt.WithExpire(time.Hour*24),
	//jwt.WithSigningKey("123456"),
	//jwt.WithSigningMethod(jwt.HS384),
	)

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
		r.Use(middleware.CircuitBreaker(
			// set http code for circuit breaker, default already includes 500 and 503
			middleware.WithValidCode(errcode.InternalServerError.Code()),
			middleware.WithValidCode(errcode.ServiceUnavailable.Code()),
		))
	}

	// trace middleware
	if config.Get().App.EnableTrace {
		r.Use(middleware.Tracing(config.Get().App.Name))
	}

	// profile performance analysis
	if config.Get().App.EnableHTTPProfile {
		prof.Register(r, prof.WithIOWaitTime())
	}

	// validator
	binding.Validator = validator.Init()

	r.GET("/health", handlerfunc.CheckHealth)
	r.GET("/ping", handlerfunc.Ping)
	r.GET("/codes", handlerfunc.ListCodes)

	if config.Get().App.Env != "prod" {
		r.GET("/config", gin.WrapF(errcode.ShowConfig([]byte(config.Show()))))
		// access path /apis/swagger/index.html
		swagger.CustomRouter(r, "apis", docs.ApiDocs)
	}

	c := newMiddlewareConfig()

	// set up all middlewares
	for _, fn := range allMiddlewareFns {
		fn(c)
	}

	// register all routes
	for _, fn := range allRouteFns {
		fn(r, c.groupPathMiddlewares, c.singlePathMiddlewares)
	}

	return r
}

type middlewareConfig struct {
	groupPathMiddlewares  map[string][]gin.HandlerFunc // middleware functions corresponding to route group
	singlePathMiddlewares map[string][]gin.HandlerFunc // middleware functions corresponding to a single route
}

func newMiddlewareConfig() *middlewareConfig {
	return &middlewareConfig{
		groupPathMiddlewares:  make(map[string][]gin.HandlerFunc),
		singlePathMiddlewares: make(map[string][]gin.HandlerFunc),
	}
}

func (c *middlewareConfig) setGroupPath(groupPath string, handlers ...gin.HandlerFunc) { //nolint
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

func (c *middlewareConfig) setSinglePath(method string, singlePath string, handlers ...gin.HandlerFunc) { //nolint
	if method == "" || singlePath == "" {
		return
	}

	key := getSinglePathKey(method, singlePath)
	handlerFns, ok := c.singlePathMiddlewares[key]
	if !ok {
		c.singlePathMiddlewares[key] = handlers
		return
	}

	c.singlePathMiddlewares[key] = append(handlerFns, handlers...)
}

func getSinglePathKey(method string, singlePath string) string { //nolint
	return strings.ToUpper(method) + "->" + singlePath
}
