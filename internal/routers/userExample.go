package routers

import (
	"github.com/zhufuyi/sponge/internal/handler"

	"github.com/gin-gonic/gin"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		userExampleRouter(group, handler.NewUserExampleHandler())
	})
}

func userExampleRouter(group *gin.RouterGroup, h handler.UserExampleHandler) {
	group.POST("/userExample", h.Create)
	group.DELETE("/userExample/:id", h.DeleteByID)
	group.POST("/userExamples/delete/ids", h.DeleteByIDs)
	group.PUT("/userExample/:id", h.UpdateByID)
	group.GET("/userExample/:id", h.GetByID)
	group.POST("/userExamples/ids", h.ListByIDs)
	group.POST("/userExamples", h.List)
}
