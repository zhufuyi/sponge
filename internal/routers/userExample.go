package routers

import (
	"github.com/zhufuyi/sponge/internal/handler"

	"github.com/gin-gonic/gin"
)

// nolint
func init() {
	routerFns = append(routerFns, func() {
		userExampleRouter(apiV1, handler.NewUserExampleHandler()) // 加入到路由组
	})
}

func userExampleRouter(group *gin.RouterGroup, h handler.UserExampleHandler) {
	group.POST("/userExample", h.Create)
	group.DELETE("/userExample/:id", h.DeleteByID)
	group.PUT("/userExample/:id", h.UpdateByID)
	group.GET("/userExample/:id", h.GetByID)
	group.POST("/userExamples/ids", h.ListByIDs)
	group.POST("/userExamples", h.List) // 通过post任意列组合查询
}
