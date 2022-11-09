package handler

import (
	"text/template"
)

func init() {
	var err error
	handlerTmpl, err = template.New("iHandler").Parse(handlerTmplRaw)
	if err != nil {
		panic(err)
	}
	routerTmpl, err = template.New("handlerRouter").Parse(routerTmplRaw)
	if err != nil {
		panic(err)
	}
}

var (
	pkgImportTmplRaw = `package handler

import (
	"context"

	serverNameExampleV1 "moduleNameExample/api/serverNameExample/v1"

	//"github.com/zhufuyi/sponge/pkg/gin/middleware"
)

`

	handlerTmpl    *template.Template
	handlerTmplRaw = `var _ serverNameExampleV1.{{$.Name}}Logicer = (*{{$.LowerName}}Handler)(nil)

type {{$.LowerName}}Handler struct {
	// example: 
	// 	{{$.LowerName}}Dao dao.{{$.Name}}Dao

	// If required, fill in the definition of the other dao code here.
}

// New{{$.Name}}Handler creating handler
func New{{$.Name}}Handler() serverNameExampleV1.{{$.Name}}Logicer {
	return &{{$.LowerName}}Handler{
		// example:
		// 	{{$.LowerName}}Dao: dao.New{{$.Name}}Dao(
		// 		model.GetDB(),
		// 		cache.New{{$.Name}}Cache(model.GetCacheType()),
		// 	),

		// If required, fill in the code to implement other dao here.
	}
}
`

	routerTmpl    *template.Template
	routerTmplRaw = `package routers

import (
	serverNameExampleV1 "moduleNameExample/api/serverNameExample/v1"
	"moduleNameExample/internal/handler"

	"github.com/zhufuyi/sponge/pkg/logger"

	"github.com/gin-gonic/gin"
)

func init() {
	rootRouterFns = append(rootRouterFns, func(r *gin.Engine) {
{{- range .ServiceNames}}
		{{.LowerName}}Router(r, handler.New{{.Name}}Handler())
{{- end}}
	})
}

{{- range .ServiceNames}}

func {{.LowerName}}Router(r *gin.Engine, iService serverNameExampleV1.{{.Name}}Logicer) {
	serverNameExampleV1.Register{{.Name}}Router(r, iService,
		serverNameExampleV1.With{{.Name}}RPCResponse(),
		serverNameExampleV1.With{{.Name}}Logger(logger.Get()),
	)
}
{{- end}}
`
)
