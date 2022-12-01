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

	serviceLogicTestTmpl    *template.Template
	serviceLogicTestTmplRaw = `package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	serverNameExampleV1 "moduleNameExample/api/serverNameExample/v1"
	"moduleNameExample/configs"
	"moduleNameExample/internal/config"

	"github.com/zhufuyi/sponge/pkg/consulcli"
	"github.com/zhufuyi/sponge/pkg/etcdcli"
	"github.com/zhufuyi/sponge/pkg/grpc/grpccli"
	"github.com/zhufuyi/sponge/pkg/nacoscli"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/consul"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/etcd"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/nacos"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func initServerNameExampleClient() *grpc.ClientConn {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		panic(err)
	}
	endpoint := fmt.Sprintf("127.0.0.1:%d", config.Get().Grpc.Port)

	var cliOptions = []grpccli.Option{
		grpccli.WithEnableLog(zap.NewNop()),
		//grpccli.WithEnableLoadBalance(),
		//grpccli.WithEnableRetry(),
	}
	if config.Get().App.RegistryDiscoveryType != "" {
		var iDiscovery registry.Discovery
		endpoint = "discovery:///" + config.Get().App.Name // Connecting to grpc services by service name

		// Use consul service discovery, note that the host field in the configuration file serverNameExample.yml
		// needs to be filled with the local ip, not 127.0.0.1, to do the health check
		if config.Get().App.RegistryDiscoveryType == "consul" {
			cli, err := consulcli.Init(config.Get().Consul.Addr, consulcli.WithWaitTime(time.Second*2))
			if err != nil {
				panic(err)
			}
			iDiscovery = consul.New(cli)
		}

		// Use etcd service discovery, use the command etcdctl get / --prefix to see if the service is registered before testing,
		// note: the IDE using a proxy may cause the connection to the etcd service to fail
		if config.Get().App.RegistryDiscoveryType == "etcd" {
			cli, err := etcdcli.Init(config.Get().Etcd.Addrs, etcdcli.WithDialTimeout(time.Second*2))
			if err != nil {
				panic(err)
			}
			iDiscovery = etcd.New(cli)
		}

		// Use nacos service discovery
		if config.Get().App.RegistryDiscoveryType == "nacos" {
			// example: endpoint = "discovery:///serverName.scheme"
			endpoint = "discovery:///" + config.Get().App.Name + ".grpc"
			cli, err := nacoscli.NewNamingClient(
				config.Get().NacosRd.IPAddr,
				config.Get().NacosRd.Port,
				config.Get().NacosRd.NamespaceID)
			if err != nil {
				panic(err)
			}
			iDiscovery = nacos.New(cli)
		}

		cliOptions = append(cliOptions, grpccli.WithDiscovery(iDiscovery))
	}

	if config.Get().App.EnableTracing {
		cliOptions = append(cliOptions, grpccli.WithEnableTrace())
	}
	if config.Get().App.EnableCircuitBreaker {
		cliOptions = append(cliOptions, grpccli.WithEnableCircuitBreaker())
	}
	if config.Get().App.EnableMetrics {
		cliOptions = append(cliOptions, grpccli.WithEnableMetrics())
	}

	conn, err := grpccli.DialInsecure(context.Background(), endpoint, cliOptions...)
	if err != nil {
		panic(err)
	}

	return conn
}

{{- range .PbServices}}

func Test_{{.LowerName}}_methods(t *testing.T) {
	conn := initServerNameExampleClient()
	cli := serverNameExampleV1.New{{.Name}}Client(conn)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)

	tests := []struct {
		name    string
		fn      func() (interface{}, error)
		wantErr bool
	}{
{{- range .Methods}}
		{
			name: "{{.MethodName}}",
			fn: func() (interface{}, error) {
				// todo enter parameters before testing
				req := &serverNameExampleV1.{{.MethodName}}Request{
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
			t.Logf("reply data: %+v", got)
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
