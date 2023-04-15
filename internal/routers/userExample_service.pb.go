package routers

import (
	serverNameExampleV1 "github.com/zhufuyi/sponge/api/serverNameExample/v1"
	"github.com/zhufuyi/sponge/internal/service"

	"github.com/zhufuyi/sponge/pkg/logger"

	"github.com/gin-gonic/gin"
)

func init() {
	apiV1RouterFns_pbExample = append(apiV1RouterFns_pbExample, func(prePath string, group *gin.RouterGroup) {
		userExampleServiceRouter(prePath, group, service.NewUserExampleServiceClient())
	})
}

func userExampleServiceRouter(prePath string, group *gin.RouterGroup, iService serverNameExampleV1.UserExampleServiceLogicer) {
	serverNameExampleV1.RegisterUserExampleServiceRouter(prePath, group, iService,
		serverNameExampleV1.WithUserExampleServiceRPCResponse(),
		serverNameExampleV1.WithUserExampleServiceLogger(logger.Get()),
	)
}
