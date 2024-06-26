package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/zhufuyi/sponge/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		userExampleRouter(group, handler.NewUserExampleHandler())
	})
}

func userExampleRouter(group *gin.RouterGroup, h handler.UserExampleHandler) {
	//group.Use(middleware.Auth()) // all of the following routes use jwt authentication
	// or group.Use(middleware.Auth(middleware.WithVerify(verify))) // token authentication

	group.POST("/userExample", h.Create)
	group.DELETE("/userExample/:id", h.DeleteByID)
	group.PUT("/userExample/:id", h.UpdateByID)
	group.GET("/userExample/:id", h.GetByID)
	group.POST("/userExample/list", h.List)
}
