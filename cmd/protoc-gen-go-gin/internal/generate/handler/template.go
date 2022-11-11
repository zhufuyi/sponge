package handler

import (
	"math/rand"
	"text/template"
	"time"
)

func init() {
	var err error
	handlerLogicTmpl, err = template.New("handlerLogic").Parse(handlerLogicTmplRaw)
	if err != nil {
		panic(err)
	}
	routerTmpl, err = template.New("handlerRouter").Parse(routerTmplRaw)
	if err != nil {
		panic(err)
	}
	httpErrCodeTmpl, err = template.New("httpErrCode").Parse(httpErrCodeTmplRaw)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
}

var (
	handlerLogicTmpl    *template.Template
	handlerLogicTmplRaw = `package handler

import (
	"context"

	serverNameExampleV1 "moduleNameExample/api/serverNameExample/v1"

	//"github.com/zhufuyi/sponge/pkg/gin/middleware"
)

{{- range .PbServices}}

var _ serverNameExampleV1.{{.Name}}Logicer = (*{{.LowerName}}Handler)(nil)

type {{.LowerName}}Handler struct {
	// example: 
	// 	{{.LowerName}}Dao dao.{{.Name}}Dao

	// If required, fill in the definition of the other dao code here.
}

// New{{.Name}}Handler creating handler
func New{{.Name}}Handler() serverNameExampleV1.{{.Name}}Logicer {
	return &{{.LowerName}}Handler{
		// example:
		// 	{{.LowerName}}Dao: dao.New{{.Name}}Dao(
		// 		model.GetDB(),
		// 		cache.New{{.Name}}Cache(model.GetCacheType()),
		// 	),

		// If required, fill in the code to implement other dao here.
	}
}

{{- range .Methods}}

func (h *{{.LowerServiceName}}Handler) {{.MethodName}}(ctx context.Context, req *serverNameExampleV1.{{.Request}}) (*serverNameExampleV1.{{.Reply}}, error) {
	// example:
	// 	reply, err := h.{{.LowerServiceName}}Dao.{{.MethodName}}(ctx, req)
	// 	if err != nil {
	//			logger.Warn("invoke error", logger.Err(err), middleware.CtxRequestIDField(ctx))
	//			return nil, ecode.InternalServerError.Err()
	//		}
	// 	return reply, nil
	//
	// If required, fill in the code for getting data from other dao here

	panic("implement me")
}

{{- end}}

{{- end}}
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
{{- range .PbServices}}
		{{.LowerName}}Router(r, handler.New{{.Name}}Handler())
{{- end}}
	})
}

{{- range .PbServices}}

func {{.LowerName}}Router(r *gin.Engine, iService serverNameExampleV1.{{.Name}}Logicer) {
	serverNameExampleV1.Register{{.Name}}Router(r, iService,
		serverNameExampleV1.With{{.Name}}RPCResponse(),
		serverNameExampleV1.With{{.Name}}Logger(logger.Get()),
	)
}
{{- end}}
`

	httpErrCodeTmpl    *template.Template
	httpErrCodeTmplRaw = `package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

{{- range .PbServices}}

// {{.LowerName}} http service level error code
var (
	{{.LowerName}}NO       = {{.RandNumber}} // number range 1~100, if there is the same number, trigger panic.
	{{.LowerName}}Name     = "{{.LowerName}}"
	{{.LowerName}}BaseCode = errcode.HCode({{.LowerName}}NO)
// --blank line--
{{- range $i, $v := .Methods}}
	Err{{.MethodName}}{{.ServiceName}}   = errcode.NewError({{.LowerServiceName}}BaseCode+{{$v.AddOne $i}}, "failed to {{.MethodName}} "+{{.LowerServiceName}}Name)
{{- end}}
	// add +1 to the previous error code
)

{{- end}}
`
)
