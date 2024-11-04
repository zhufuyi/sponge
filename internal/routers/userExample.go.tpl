package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/zhufuyi/sponge/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		{{.TableNameCamelFCL}}Router(group, handler.New{{.TableNameCamel}}Handler())
	})
}

func {{.TableNameCamelFCL}}Router(group *gin.RouterGroup, h handler.{{.TableNameCamel}}Handler) {
	g := group.Group("/{{.TableNameCamelFCL}}")

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithVerify(fn))
	//g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/", h.Create)          // [post] /api/v1/{{.TableNameCamelFCL}}
	g.DELETE("/:{{.ColumnNameCamelFCL}}", h.DeleteBy{{.ColumnNameCamel}}) // [delete] /api/v1/{{.TableNameCamelFCL}}/:{{.ColumnNameCamelFCL}}
	g.PUT("/:{{.ColumnNameCamelFCL}}", h.UpdateBy{{.ColumnNameCamel}})    // [put] /api/v1/{{.TableNameCamelFCL}}/:{{.ColumnNameCamelFCL}}
	g.GET("/:{{.ColumnNameCamelFCL}}", h.GetBy{{.ColumnNameCamel}})       // [get] /api/v1/{{.TableNameCamelFCL}}/:{{.ColumnNameCamelFCL}}
	g.POST("/list", h.List)        // [post] /api/v1/{{.TableNameCamelFCL}}/list
}
