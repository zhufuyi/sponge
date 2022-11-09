package routers

import (
	serverNameExampleV1 "github.com/zhufuyi/sponge/api/serverNameExample/v1"
	"github.com/zhufuyi/sponge/internal/service"

	"github.com/zhufuyi/sponge/pkg/logger"

	"github.com/gin-gonic/gin"
)

func init() {
	rootRouterFns = append(rootRouterFns, func(r *gin.Engine) {
		userExampleServiceRouter(r, service.NewUserExampleServiceClient())
	})
}

func userExampleServiceRouter(r *gin.Engine, iService serverNameExampleV1.UserExampleServiceLogicer) {
	serverNameExampleV1.RegisterUserExampleServiceRouter(r, iService,
		serverNameExampleV1.WithUserExampleServiceRPCResponse(),
		serverNameExampleV1.WithUserExampleServiceLogger(logger.Get()),
	)
}
