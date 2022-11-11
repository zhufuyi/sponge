package service

import (
	"math/rand"
	"text/template"
	"time"
)

func init() {
	var err error
	serviceLogicTmpl, err = template.New("serviceLogic").Parse(serviceLogicTmplRaw)
	if err != nil {
		panic(err)
	}
	routerTmpl, err = template.New("serviceRouter").Parse(routerTmplRaw)
	if err != nil {
		panic(err)
	}
	rpcErrCodeTmpl, err = template.New("httpErrCode").Parse(rpcErrCodeTmplRaw)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
}

var (
	serviceLogicTmpl    *template.Template
	serviceLogicTmplRaw = `package service

import (
	"context"

	serverNameExampleV1 "moduleNameExample/api/serverNameExample/v1"
	"moduleNameExample/internal/rpcclient"
)

{{- range .PbServices}}

var _ serverNameExampleV1.{{.Name}}Logicer = (*{{.LowerName}}Client)(nil)

type {{.LowerName}}Client struct {
	{{.LowerName}}Cli serverNameExampleV1.{{.Name}}Client
	// If required, fill in the definition of the other service client code here.
}

// New{{.Name}}Client creating rpc clients
func New{{.Name}}Client() serverNameExampleV1.{{.Name}}Logicer {
	return &{{.LowerName}}Client{
		{{.LowerName}}Cli: serverNameExampleV1.New{{.Name}}Client(rpcclient.GetServerNameExampleRPCConn()),
		// If required, fill in the code to implement other service clients here.
	}
}

{{- range .Methods}}

func (c *{{.LowerServiceName}}Client) {{.MethodName}}(ctx context.Context, req *serverNameExampleV1.{{.Request}}) (*serverNameExampleV1.{{.Reply}}, error) {
	// implement me
	// If required, fill in the code to fetch data from other rpc servers here.

	return c.{{.LowerServiceName}}Cli.{{.MethodName}}(ctx, req)
}

{{- end}}

{{- end}}
`

	routerTmpl    *template.Template
	routerTmplRaw = `package routers

import (
	serverNameExampleV1 "moduleNameExample/api/serverNameExample/v1"
	"moduleNameExample/internal/service"

	"github.com/zhufuyi/sponge/pkg/logger"

	"github.com/gin-gonic/gin"
)

func init() {
	rootRouterFns = append(rootRouterFns, func(r *gin.Engine) {
{{- range .PbServices}}
		{{.LowerName}}Router(r, service.New{{.Name}}Client())
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

	rpcErrCodeTmpl    *template.Template
	rpcErrCodeTmplRaw = `package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

{{- range .PbServices}}

// {{.LowerName}} rpc service level error code
var (
	_{{.LowerName}}NO       = {{.RandNumber}} // number range 1~100, if there is the same number, trigger panic.
	_{{.LowerName}}Name     = "{{.LowerName}}"
	_{{.LowerName}}BaseCode = errcode.HCode(_{{.LowerName}}NO)
// --blank line--
{{- range $i, $v := .Methods}}
	Status{{.MethodName}}{{.ServiceName}}   = errcode.NewError(_{{.LowerServiceName}}BaseCode+{{$v.AddOne $i}}, "failed to {{.MethodName}} "+_{{.LowerServiceName}}Name)
{{- end}}
	// add +1 to the previous error code
)

{{- end}}
`
)
