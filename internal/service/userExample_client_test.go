// 开启grpc服务端后再进行测试，下面对userExample各个方法进行测试和
// 压测(复制压测报告文件路径到浏览器查看)

package service

import (
	"context"
	"fmt"
	"github.com/zhufuyi/sponge/pkg/etcdcli"
	"github.com/zhufuyi/sponge/pkg/nacoscli"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/etcd"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/nacos"
	"testing"
	"time"

	pb "github.com/zhufuyi/sponge/api/serverNameExample/v1"
	"github.com/zhufuyi/sponge/api/types"
	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/config"

	"github.com/zhufuyi/sponge/pkg/consulcli"
	"github.com/zhufuyi/sponge/pkg/grpc/benchmark"
	"github.com/zhufuyi/sponge/pkg/grpc/grpccli"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/consul"
	"go.uber.org/zap"
)

func initUserExampleServiceClient() pb.UserExampleServiceClient {
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
	if config.Get().App.EnableRegistryDiscovery {
		var iDiscovery registry.Discovery
		endpoint = "discovery:///" + config.Get().App.Name // 通过服务名称连接grpc服务

		// 使用consul做发现，注意配置文件serverNameExample.yml的字段host需要填本机ip，不是127.0.0.1，用来做健康检查
		if config.Get().App.RegistryDiscoveryType == "consul" {
			cli, err := consulcli.Init(config.Get().Consul.Addr, consulcli.WithWaitTime(time.Second*2))
			if err != nil {
				panic(err)
			}
			iDiscovery = consul.New(cli)
		}

		// 使用etcd做服务发现，测试前使用命令etcdctl get / --prefix查看是否有服务注册，注：IDE使用代理可能会造成连接etcd服务失败
		if config.Get().App.RegistryDiscoveryType == "etcd" {
			cli, err := etcdcli.Init(config.Get().Etcd.Addrs, etcdcli.WithDialTimeout(time.Second*2))
			if err != nil {
				panic(err)
			}
			iDiscovery = etcd.New(cli)
		}

		// 使用nacos做服务发现
		if config.Get().App.RegistryDiscoveryType == "nacos" {
			// example: endpoint = "discovery:///serverName.scheme"
			endpoint += ".grpc" // 在endpoint+.grpc
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
	//defer conn.Close()

	return pb.NewUserExampleServiceClient(conn)
}

// 通过客户端测试userExample的各个方法
func Test_userExampleService_methods(t *testing.T) {
	cli := initUserExampleServiceClient()
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)

	tests := []struct {
		name    string
		fn      func() (interface{}, error)
		wantErr bool
	}{
		// todo generate the service struct code here
		// delete the templates code start
		{
			name: "Create",
			fn: func() (interface{}, error) {
				// todo test after filling in parameters
				return cli.Create(ctx, &pb.CreateUserExampleRequest{
					Name:     "foo7",
					Email:    "foo7@bar.com",
					Password: "f447b20a7fcbf53a5d5be013ea0b15af",
					Phone:    "16000000000",
					Avatar:   "http://internal.com/7.jpg",
					Age:      11,
					Gender:   2,
				})
			},
			wantErr: false,
		},

		{
			name: "UpdateByID",
			fn: func() (interface{}, error) {
				// todo test after filling in parameters
				return cli.UpdateByID(ctx, &pb.UpdateUserExampleByIDRequest{
					Id:    7,
					Phone: "16000000001",
					Age:   11,
				})
			},
			wantErr: false,
		},
		// delete the templates code end
		{
			name: "DeleteByID",
			fn: func() (interface{}, error) {
				// todo test after filling in parameters
				return cli.DeleteByID(ctx, &pb.DeleteUserExampleByIDRequest{
					Id: 100,
				})
			},
			wantErr: false,
		},

		{
			name: "GetByID",
			fn: func() (interface{}, error) {
				// todo test after filling in parameters
				return cli.GetByID(ctx, &pb.GetUserExampleByIDRequest{
					Id: 1,
				})
			},
			wantErr: false,
		},

		{
			name: "ListByIDs",
			fn: func() (interface{}, error) {
				// todo test after filling in parameters
				return cli.ListByIDs(ctx, &pb.ListUserExampleByIDsRequest{
					Ids: []uint64{1, 2, 3},
				})
			},
			wantErr: false,
		},

		{
			name: "List",
			fn: func() (interface{}, error) {
				// todo test after filling in parameters
				return cli.List(ctx, &pb.ListUserExampleRequest{
					Params: &types.Params{
						Page:  0,
						Limit: 10,
						Sort:  "",
						Columns: []*types.Column{
							{
								Name:  "id",
								Exp:   ">=",
								Value: "1",
								Logic: "",
							},
						},
					},
				})
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fn()
			if (err != nil) != tt.wantErr {
				// 如果没有开启rpc服务端，会报错transport: Error while dialing dial tcp......，这里忽略测试错误
				t.Logf("test '%s' error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			t.Logf("reply data: %+v", got)
		})
	}
}

// 压测userExample的各个方法，完成后复制报告路径到浏览器查看
func Test_userExampleService_benchmark(t *testing.T) {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		panic(err)
	}
	host := fmt.Sprintf("127.0.0.1:%d", config.Get().Grpc.Port)
	protoFile := configs.Path("../api/serverNameExample/v1/userExample.proto")
	// 如果压测过程中缺少第三方依赖，复制到项目的third_party目录下(不包括import路径)
	importPaths := []string{
		configs.Path("../third_party"), // third_party目录
		configs.Path(".."),             // third_party的上一级目录
	}

	tests := []struct {
		name    string
		fn      func() error
		wantErr bool
	}{
		{
			name: "GetByID",
			fn: func() error {
				// todo test after filling in parameters
				message := &pb.GetUserExampleByIDRequest{
					Id: 1,
				}
				b, err := benchmark.New(host, protoFile, "GetByID", message, 1000, importPaths...)
				if err != nil {
					return err
				}
				return b.Run()
			},
			wantErr: false,
		},

		{
			name: "ListByIDs",
			fn: func() error {
				// todo test after filling in parameters
				message := &pb.ListUserExampleByIDsRequest{
					Ids: []uint64{1, 2, 3},
				}
				b, err := benchmark.New(host, protoFile, "ListByIDs", message, 1000, importPaths...)
				if err != nil {
					return err
				}
				return b.Run()
			},
			wantErr: false,
		},

		{
			name: "List",
			fn: func() error {
				// todo test after filling in parameters
				message := &pb.ListUserExampleRequest{
					Params: &types.Params{
						Page:  0,
						Limit: 10,
						Sort:  "",
						Columns: []*types.Column{
							{
								Name:  "id",
								Exp:   ">=",
								Value: "1",
								Logic: "",
							},
						},
					},
				}
				b, err := benchmark.New(host, protoFile, "List", message, 100, importPaths...)
				if err != nil {
					return err
				}
				return b.Run()
			},
			wantErr: false,
		},
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
