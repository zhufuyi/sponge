package service

import (
	"text/template"
)

func init() {
	var err error
	serviceTmpl, err = template.New("iService").Parse(serviceTmplRaw)
	if err != nil {
		panic(err)
	}
	routerTmpl, err = template.New("serviceRouter").Parse(routerTmplRaw)
	if err != nil {
		panic(err)
	}
}

var (
	pkgImportTmplRaw = `package service

import (
	"context"

	serverNameExampleV1 "moduleNameExample/api/serverNameExample/v1"
	"moduleNameExample/internal/rpcclient"
)

`

	serviceTmpl    *template.Template
	serviceTmplRaw = `var _ serverNameExampleV1.{{$.Name}}Logicer = (*{{$.LowerName}}Client)(nil)

type {{$.LowerName}}Client struct {
	{{$.LowerName}}Cli serverNameExampleV1.{{$.Name}}Client
	// If required, fill in the definition of the other service client code here.
}

// New{{$.Name}}Client creating rpc clients
func New{{$.Name}}Client() serverNameExampleV1.{{$.Name}}Logicer {
	return &{{$.LowerName}}Client{
		{{$.LowerName}}Cli: serverNameExampleV1.New{{$.Name}}Client(rpcclient.GetServerNameExampleRPCConn()),
		// If required, fill in the code to implement other service clients here.
	}
}
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
{{- range .ServiceNames}}
		{{.LowerName}}Router(r, service.New{{.Name}}Client())
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
