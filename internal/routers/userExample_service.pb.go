package routers

import (
	"context"

	serverNameExampleV1 "github.com/zhufuyi/sponge/api/serverNameExample/v1"
	"github.com/zhufuyi/sponge/internal/service"

	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/logger"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

func init() {
	allMiddlewareFns = append(allMiddlewareFns, func(c *middlewareConfig) {
		userExampleMiddlewares(c)
	})

	allRouteFns = append(allRouteFns,
		func(r *gin.Engine, groupPathMiddlewares map[string][]gin.HandlerFunc, singlePathMiddlewares map[string][]gin.HandlerFunc) {
			userExampleServiceRouter(r, groupPathMiddlewares, singlePathMiddlewares, service.NewUserExampleServiceClient())
		})
}

func userExampleServiceRouter(
	r *gin.Engine,
	groupPathMiddlewares map[string][]gin.HandlerFunc,
	singlePathMiddlewares map[string][]gin.HandlerFunc,
	iService serverNameExampleV1.UserExampleServiceLogicer) {
	fn := func(c *gin.Context) context.Context {
		md := metadata.New(map[string]string{
			// set metadata to be passed from http to rpc
			middleware.ContextRequestIDKey: middleware.GCtxRequestID(c), // request_id
			//middleware.HeaderAuthorizationKey: c.GetHeader(middleware.HeaderAuthorizationKey),  // authorization
		})
		return metadata.NewOutgoingContext(c, md)
	}

	serverNameExampleV1.RegisterUserExampleServiceRouter(
		r,
		groupPathMiddlewares,
		singlePathMiddlewares,
		iService,
		serverNameExampleV1.WithUserExampleServiceRPCResponse(),
		serverNameExampleV1.WithUserExampleServiceLogger(logger.Get()),
		serverNameExampleV1.WithUserExampleServiceRPCStatusToHTTPCode(
		// Set some error codes to standard http return codes,
		// by default there is already ecode.StatusInternalServerError and ecode.StatusServiceUnavailable
		// example:
		// 	ecode.StatusUnimplemented, ecode.StatusAborted,
		),
		serverNameExampleV1.WithUserExampleServiceWrapCtx(fn),
	)
}

// you can set the middleware of a route group, or set the middleware of a single route,
// or you can mix them, pay attention to the duplication of middleware when mixing them,
// it is recommended to set the middleware of a single route in preference
func userExampleMiddlewares(c *middlewareConfig) {
	// set up group route middleware, group path is left prefix rules,
	// if the left prefix is hit, the middleware will take effect, e.g. group route /api/v1, route /api/v1/userExample/:id  will take effect
	// c.setGroupPath("/api/v1/userExample", middleware.Auth())

	// set up single route middleware, just uncomment the code and fill in the middlewares, nothing else needs to be changed
	//c.setSinglePath("GET", "/api/v1/userExample/:id", middleware.Auth())
	//c.setSinglePath("DELETE", "/api/v1/userExample/:id", middleware.Auth())
	//c.setSinglePath("PUT", "/api/v1/userExample/:id", middleware.Auth())
}
