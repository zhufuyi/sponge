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
	serviceLogicTestTmpl, err = template.New("serviceLogicTestTmpl").Parse(serviceLogicTestTmplRaw)
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

	//"github.com/zhufuyi/sponge/pkg/grpc/interceptor"
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
	//	    err := req.Validate()
	//	    if err != nil {
	//		    logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
	//		    return nil, ecode.StatusInvalidParams.Err()
	//	    }
    //
	// 	reply, err := s.xxxDao.XxxMethod(ctx, req)
	// 	if err != nil {
	//			logger.Warn("XxxMethod error", logger.Err(err), interceptor.ServerCtxRequestIDField(ctx))
	//			return nil, ecode.InternalServerError.Err()
	//		}
	// 	return reply, nil

	// fill in the business logic code

	panic("implement me")
}

{{- end}}

{{- end}}
`

	serviceLogicTestTmpl    *template.Template
	serviceLogicTestTmplRaw = `package service

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	serverNameExampleV1 "moduleNameExample/api/serverNameExample/v1"
	"moduleNameExample/configs"
	"moduleNameExample/internal/config"

	"github.com/zhufuyi/sponge/pkg/grpc/benchmark"
)

{{- range .PbServices}}

// Test each method of {{.LowerName}} via the rpc client
func Test_service_{{.LowerName}}_methods(t *testing.T) {
	conn := getRPCClientConnForTest()
	cli := serverNameExampleV1.New{{.Name}}Client(conn)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	tests := []struct {
		name    string
		fn      func() (interface{}, error)
		wantErr bool
	}{
{{- range .Methods}}
		{
			name: "{{.MethodName}}",
			fn: func() (interface{}, error) {
				// todo type in the parameters to test
				req := &serverNameExampleV1.{{.Request}}{
{{- range .RequestFields}}
					{{.FieldName}}: {{.GoTypeZero}}, {{if .Comment}} {{.Comment}}{{end}}
{{- end}}
				}
				return cli.{{.MethodName}}(ctx, req)
			},
			wantErr: false,
		},
{{- end}}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fn()
			if (err != nil) != tt.wantErr {
				t.Errorf("test '%s' error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			data, _ := json.MarshalIndent(got, "", "    ")
			fmt.Println(string(data))
		})
	}
}

// Perform a stress test on {{.LowerName}}'s method and 
// copy the press test report to your browser when you are finished.
func Test_service_{{.LowerName}}_benchmark(t *testing.T) {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		panic(err)
	}
	host := fmt.Sprintf("127.0.0.1:%d", config.Get().Grpc.Port)
	protoFile := configs.Path("../api/serverNameExample/v1/{{.ProtoName}}")
	// If third-party dependencies are missing during the press test,
	// copy them to the project's third_party directory.
	importPaths := []string{
		configs.Path("../third_party"), // third_party directory
		configs.Path(".."),             // Previous level of third_party
	}

	tests := []struct {
		name    string
		fn      func() error
		wantErr bool
	}{
{{- range .Methods}}
		{
			name: "{{.MethodName}}",
			fn: func() error {
				// todo type in the parameters to test
				message := &serverNameExampleV1.{{.Request}}{
{{- range .RequestFields}}
					{{.FieldName}}: {{.GoTypeZero}}, {{if .Comment}} {{.Comment}}{{end}}
{{- end}}
				}
				var total uint = 1000 // total number of requests
				b, err := benchmark.New(host, protoFile, "{{.MethodName}}", message, total, importPaths...)
				if err != nil {
					return err
				}
				return b.Run()
			},
			wantErr: false,
		},
{{- end}}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if (err != nil) != tt.wantErr {
				t.Errorf("test '%s' error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
		})
	}
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
