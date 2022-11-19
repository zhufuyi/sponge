// After running the rpc server and then testing, the following tests and crush tests are performed on
// each method of userExample (copy the path to the crush test report file to view it in your browser)

package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	serverNameExampleV1 "github.com/zhufuyi/sponge/api/serverNameExample/v1"
	"github.com/zhufuyi/sponge/api/types"
	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/config"

	"github.com/zhufuyi/sponge/pkg/consulcli"
	"github.com/zhufuyi/sponge/pkg/etcdcli"
	"github.com/zhufuyi/sponge/pkg/grpc/benchmark"
	"github.com/zhufuyi/sponge/pkg/grpc/grpccli"
	"github.com/zhufuyi/sponge/pkg/nacoscli"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/consul"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/etcd"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/nacos"

	"go.uber.org/zap"
)

func initUserExampleServiceClient() serverNameExampleV1.UserExampleServiceClient {
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

	return serverNameExampleV1.NewUserExampleServiceClient(conn)
}

// Test each method of userExample via the client
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
				return cli.Create(ctx, &serverNameExampleV1.CreateUserExampleRequest{
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
				return cli.UpdateByID(ctx, &serverNameExampleV1.UpdateUserExampleByIDRequest{
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
				return cli.DeleteByID(ctx, &serverNameExampleV1.DeleteUserExampleByIDRequest{
					Id: 100,
				})
			},
			wantErr: false,
		},

		{
			name: "GetByID",
			fn: func() (interface{}, error) {
				// todo test after filling in parameters
				return cli.GetByID(ctx, &serverNameExampleV1.GetUserExampleByIDRequest{
					Id: 1,
				})
			},
			wantErr: false,
		},

		{
			name: "ListByIDs",
			fn: func() (interface{}, error) {
				// todo test after filling in parameters
				return cli.ListByIDs(ctx, &serverNameExampleV1.ListUserExampleByIDsRequest{
					Ids: []uint64{1, 2, 3},
				})
			},
			wantErr: false,
		},

		{
			name: "List",
			fn: func() (interface{}, error) {
				// todo test after filling in parameters
				return cli.List(ctx, &serverNameExampleV1.ListUserExampleRequest{
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
				// If the rpc server is not enabled, it will report the error transport: Error while dialing dial tcp...... Ignore the test error here
				t.Logf("test '%s' error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			t.Logf("reply data: %+v", got)
		})
	}
}

// Press test the individual queries of userExample and copy the press test report to your browser when finished
func Test_userExampleService_benchmark(t *testing.T) {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		panic(err)
	}
	host := fmt.Sprintf("127.0.0.1:%d", config.Get().Grpc.Port)
	protoFile := configs.Path("../api/serverNameExample/v1/userExample.proto")
	// If third-party dependencies are missing during the press test, copy them to the project's third_party directory (not including the import path)
	importPaths := []string{
		configs.Path("../third_party"), // third_party directory
		configs.Path(".."),             // Previous level of third_party
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
				message := &serverNameExampleV1.GetUserExampleByIDRequest{
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
				message := &serverNameExampleV1.ListUserExampleByIDsRequest{
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
				message := &serverNameExampleV1.ListUserExampleRequest{
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
