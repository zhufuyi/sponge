package routers

import (
	"context"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"

	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/logger"

	serverNameExampleV1 "github.com/go-dev-frame/sponge/api/serverNameExample/v1"
	"github.com/go-dev-frame/sponge/internal/service"
)

func init() {
	allMiddlewareFns = append(allMiddlewareFns, func(c *middlewareConfig) {
		userExampleMiddlewares(c)
	})

	allRouteFns = append(allRouteFns,
		func(r *gin.Engine, groupPathMiddlewares map[string][]gin.HandlerFunc, singlePathMiddlewares map[string][]gin.HandlerFunc) {
			userExampleServiceRouter(r, groupPathMiddlewares, singlePathMiddlewares, service.NewUserExampleClient())
		})
}

func userExampleServiceRouter(
	r *gin.Engine,
	groupPathMiddlewares map[string][]gin.HandlerFunc,
	singlePathMiddlewares map[string][]gin.HandlerFunc,
	iService serverNameExampleV1.UserExampleLogicer) {
	fn := func(c *gin.Context) context.Context {
		md := metadata.New(map[string]string{
			// set metadata to be passed from http to rpc
			middleware.ContextRequestIDKey: middleware.GCtxRequestID(c), // request_id
			//middleware.HeaderAuthorizationKey: c.GetHeader(middleware.HeaderAuthorizationKey),  // authorization
		})
		return metadata.NewOutgoingContext(c.Request.Context(), md)
	}

	serverNameExampleV1.RegisterUserExampleRouter(
		r,
		groupPathMiddlewares,
		singlePathMiddlewares,
		iService,
		serverNameExampleV1.WithUserExampleLogger(logger.Get()),
		serverNameExampleV1.WithUserExampleRPCResponse(),
		serverNameExampleV1.WithUserExampleWrapCtx(fn),
		serverNameExampleV1.WithUserExampleRPCStatusToHTTPCode(
		// Set some error codes to standard http return codes,
		// by default there is already ecode.StatusInternalServerError and ecode.StatusServiceUnavailable
		// example:
		// 	ecode.StatusUnimplemented, ecode.StatusAborted,
		),
	)
}

// you can set the middleware of a route group, or set the middleware of a single route,
// or you can mix them, pay attention to the duplication of middleware when mixing them,
// it is recommended to set the middleware of a single route in preference
func userExampleMiddlewares(c *middlewareConfig) {
	_ = c

	// set up group route middleware, group path is left prefix rules,
	// if the left prefix is hit, the middleware will take effect, e.g. group route /api/v1, route /api/v1/userExample/:id  will take effect
	// c.setGroupPath("/api/v1/userExample", middleware.Auth())

	// set up single route middleware, just uncomment the code and fill in the middlewares, nothing else needs to be changed
	//c.setSinglePath("POST", "/api/v1/userExample", middleware.Auth())
	//c.setSinglePath("DELETE", "/api/v1/userExample/:id", middleware.Auth())
	//c.setSinglePath("PUT", "/api/v1/userExample/:id", middleware.Auth())
	//c.setSinglePath("GET", "/api/v1/userExample/:id", middleware.Auth())
	//c.setSinglePath("POST", "/api/v1/userExample/list", middleware.Auth())
}
