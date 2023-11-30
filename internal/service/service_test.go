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
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/nacoscli"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/consul"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/etcd"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/nacos"
	"github.com/zhufuyi/sponge/pkg/utils"

	"google.golang.org/grpc"
)

func TestRegisterAllService(t *testing.T) {
	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		server := grpc.NewServer()
		RegisterAllService(server)
		cancel()
	})
}

// The default is to connect to the local grpc service, if you want to connect to a remote grpc service,
// pass in the parameter grpcClient.
func getRPCClientConnForTest(grpcClient ...config.GrpcClient) *grpc.ClientConn {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		panic(err)
	}
	var grpcClientCfg config.GrpcClient

	if len(grpcClient) == 0 {
		grpcClientCfg = config.GrpcClient{
			Host: "127.0.0.1",
			Port: config.Get().Grpc.Port,

			Name:                  "",
			EnableLoadBalance:     false,
			RegistryDiscoveryType: "",
			ClientSecure:          config.ClientSecure{},
			ClientToken:           config.ClientToken{},
		}
	} else {
		grpcClientCfg = grpcClient[0]
	}

	endpoint := grpcClientCfg.Host + ":" + utils.IntToStr(grpcClientCfg.Port)
	var cliOptions []grpccli.Option

	// load balance
	if grpcClientCfg.EnableLoadBalance {
		cliOptions = append(cliOptions, grpccli.WithEnableLoadBalance())
	}

	// secure
	cliOptions = append(cliOptions, grpccli.WithSecure(
		grpcClientCfg.ClientSecure.Type,
		grpcClientCfg.ClientSecure.ServerName,
		grpcClientCfg.ClientSecure.CaFile,
		grpcClientCfg.ClientSecure.CertFile,
		grpcClientCfg.ClientSecure.KeyFile,
	))

	// token
	cliOptions = append(cliOptions, grpccli.WithToken(
		grpcClientCfg.ClientToken.Enable,
		grpcClientCfg.ClientToken.AppID,
		grpcClientCfg.ClientToken.AppKey,
	))

	cliOptions = append(cliOptions,
		grpccli.WithEnableRequestID(),
		grpccli.WithEnableLog(logger.Get()),
	)

	isUseDiscover := false
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
			isUseDiscover = true
		}

		// Use etcd service discovery, use the command etcdctl get / --prefix to see if the service is registered before testing,
		// note: the IDE using a proxy may cause the connection to the etcd service to fail
		if config.Get().App.RegistryDiscoveryType == "etcd" {
			cli, err := etcdcli.Init(config.Get().Etcd.Addrs, etcdcli.WithDialTimeout(time.Second*2))
			if err != nil {
				panic(err)
			}
			iDiscovery = etcd.New(cli)
			isUseDiscover = true
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
			isUseDiscover = true
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

	msg := "dialing grpc server"
	if isUseDiscover {
		msg += " with discovery from " + config.Get().App.RegistryDiscoveryType
	}
	logger.Info(msg, logger.String("name", config.Get().App.Name), logger.String("endpoint", endpoint))

	conn, err := grpccli.Dial(context.Background(), endpoint, cliOptions...)
	if err != nil {
		panic(err)
	}

	return conn
}
