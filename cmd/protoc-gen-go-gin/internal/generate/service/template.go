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
}

var (
	serviceTmpl    *template.Template
	serviceTmplRaw = `package service

import (
	"context"

	serverNameExampleV1 "module_name_example/api/server_name_example/v1"
	"module_name_example/internal/rpcclient"
)

var _ serverNameExampleV1.{{$.Name}}Logicer = (*{{$.LowerName}}Client)(nil)

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
)
