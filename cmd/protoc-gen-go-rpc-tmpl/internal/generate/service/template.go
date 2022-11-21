package service

import (
	"math/rand"
	"text/template"
	"time"
)

func init() {
	var err error
	serviceLogicTmpl, err = template.New("serviceLogicTmpl").Parse(serviceLogicTmplRaw)
	if err != nil {
		panic(err)
	}
	rpcErrCodeTmpl, err = template.New("rpcErrCode").Parse(rpcErrCodeTmplRaw)
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

	//"moduleNameExample/internal/cache"
	//"moduleNameExample/internal/dao"
	//"moduleNameExample/internal/ecode"
	//"moduleNameExample/internal/model"

	//"github.com/zhufuyi/sponge/pkg/logger"

	"google.golang.org/grpc"
)

func init() {
	registerFns = append(registerFns, func(server *grpc.Server) {
{{- range .PbServices}}
		serverNameExampleV1.Register{{.Name}}Server(server, New{{.Name}}Server())
{{- end}}
	})
}

{{- range .PbServices}}

var _ serverNameExampleV1.{{.Name}}Server = (*{{.LowerName}})(nil)

type {{.LowerName}} struct {
	serverNameExampleV1.Unimplemented{{.Name}}Server

	// example:
	//	iDao dao.{{.Name}}Dao
}

// New{{.Name}}Server create a server
func New{{.Name}}Server() serverNameExampleV1.{{.Name}}Server {
	return &{{.LowerName}}{
		// example:
		//	iDao: dao.New{{.Name}}Dao(
		//		model.GetDB(),
		//		cache.New{{.Name}}Cache(model.GetCacheType()),
		//	),
	}
}

{{- range .Methods}}

func (s *{{.LowerServiceName}}) {{.MethodName}}(ctx context.Context, req *serverNameExampleV1.{{.Request}}) (*serverNameExampleV1.{{.Reply}}, error) {
	// example:
	//	err := req.Validate()
	//	if err != nil {
	//		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req))
	//		return nil, ecode.StatusInvalidParams.Err()
	//	}

	// fill in the business code

	panic("implement me")
}

{{- end}}

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
