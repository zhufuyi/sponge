package routers

import (
	serverNameExampleV1 "github.com/zhufuyi/sponge/api/serverNameExample/v1"
	"github.com/zhufuyi/sponge/internal/service"

	"github.com/zhufuyi/sponge/pkg/logger"

	"github.com/gin-gonic/gin"
)

func init() {
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
	)
}
