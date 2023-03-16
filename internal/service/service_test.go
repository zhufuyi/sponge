package service

import (
	"context"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/config"

	"github.com/zhufuyi/sponge/pkg/consulcli"
	"github.com/zhufuyi/sponge/pkg/etcdcli"
	"github.com/zhufuyi/sponge/pkg/grpc/grpccli"
	"github.com/zhufuyi/sponge/pkg/nacoscli"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/consul"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/etcd"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/nacos"
	"github.com/zhufuyi/sponge/pkg/utils"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func TestRegisterAllService(t *testing.T) {
	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		server := grpc.NewServer()
		RegisterAllService(server)
		cancel()
	})
}

func getRPCClientConnForTest() *grpc.ClientConn {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		panic(err)
	}

	if len(config.Get().GrpcClient) == 0 {
		panic("no rpc client configuration information found in the configuration file")
	}

	// Change the configuration information of the 0th client before testing
	rpcClientCfg := config.Get().GrpcClient[0]

	endpoint := rpcClientCfg.Host + ":" + utils.IntToStr(rpcClientCfg.Port)
	var cliOptions []grpccli.Option

	// secure
	cliOptions = append(cliOptions, grpccli.WithSecure(
		rpcClientCfg.ClientSecure.Type,
		rpcClientCfg.ClientSecure.ServerName,
		rpcClientCfg.ClientSecure.CaFile,
		rpcClientCfg.ClientSecure.CertFile,
		rpcClientCfg.ClientSecure.KeyFile,
	))

	// token
	cliOptions = append(cliOptions, grpccli.WithToken(
		rpcClientCfg.ClientToken.Enable,
		rpcClientCfg.ClientToken.AppID,
		rpcClientCfg.ClientToken.AppKey,
	))

	cliOptions = append(cliOptions,
		grpccli.WithEnableRequestID(),
		grpccli.WithEnableLog(zap.NewNop()),
	)

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

	if config.Get().App.EnableTrace {
		cliOptions = append(cliOptions, grpccli.WithEnableTrace())
	}
	if config.Get().App.EnableCircuitBreaker {
		cliOptions = append(cliOptions, grpccli.WithEnableCircuitBreaker())
	}
	if config.Get().App.EnableMetrics {
		cliOptions = append(cliOptions, grpccli.WithEnableMetrics())
	}

	conn, err := grpccli.Dial(context.Background(), endpoint, cliOptions...)
	if err != nil {
		panic(err)
	}

	return conn
}
