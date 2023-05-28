package routers

import (
	serverNameExampleV1 "github.com/zhufuyi/sponge/api/serverNameExample/v1"
	"github.com/zhufuyi/sponge/internal/service"

	"github.com/zhufuyi/sponge/pkg/logger"

	"github.com/gin-gonic/gin"
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
	serverNameExampleV1.RegisterUserExampleServiceRouter(
		r,
		groupPathMiddlewares,
		singlePathMiddlewares,
		iService,
		serverNameExampleV1.WithUserExampleServiceRPCResponse(),
		serverNameExampleV1.WithUserExampleServiceLogger(logger.Get()),
		serverNameExampleV1.WithUserExampleServiceRPCStatusToHTTPCode(
		//ecode.StatusUnimplemented, ecode.StatusAborted,
		),
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
