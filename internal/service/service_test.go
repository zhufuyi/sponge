package service

import (
	"context"
	"io"
	"strconv"
	"testing"
	"time"

	"google.golang.org/grpc"

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

	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/config"
)

var ioEOF = io.EOF

func TestRegisterAllService(t *testing.T) {
	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		server := grpc.NewServer()
		RegisterAllService(server)
		cancel()
	})
}

// The default is to connect to the local grpc server, if you want to connect to a remote grpc server,
// pass in the parameter grpcClient.
func getRPCClientConnForTest(grpcClient ...config.GrpcClient) *grpc.ClientConn {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		panic(err)
	}
	var grpcClientCfg config.GrpcClient

	// custom config
	if len(grpcClient) > 0 {
		// parameter config, highest priority
		grpcClientCfg = grpcClient[0]
	} else {
		// grpcClient config in the yml file, second priority
		if len(config.Get().GrpcClient) > 0 {
			for _, v := range config.Get().GrpcClient {
				if v.Name == config.Get().App.Name { // match the current app name
					grpcClientCfg = v
					break
				}
			}
		}
	}
	// if no custom config found, use the default config
	if grpcClientCfg.Name == "" {
		grpcClientCfg = config.GrpcClient{
			Host: config.Get().App.Host,
			Port: config.Get().Grpc.Port,
			// If RegistryDiscoveryType is not empty, service discovery is used, and Host and Port values are invalid
			RegistryDiscoveryType: config.Get().App.RegistryDiscoveryType, // supports consul, etcd and nacos
			Name:                  config.Get().App.Name,
		}
		if grpcClientCfg.RegistryDiscoveryType != "" {
			grpcClientCfg.EnableLoadBalance = true
		}
	}

	var cliOptions []grpccli.Option

	if grpcClientCfg.Timeout > 0 {
		cliOptions = append(cliOptions, grpccli.WithTimeout(time.Second*time.Duration(grpcClientCfg.Timeout)))
	}

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

	var (
		endpoint      string
		isUseDiscover bool
		iDiscovery    registry.Discovery
	)

	switch grpcClientCfg.RegistryDiscoveryType {
	case "consul":
		endpoint = "discovery:///" + grpcClientCfg.Name // Connecting to grpc services by service name
		cli, err := consulcli.Init(config.Get().Consul.Addr, consulcli.WithWaitTime(time.Second*2))
		if err != nil {
			panic(err)
		}
		iDiscovery = consul.New(cli)
		isUseDiscover = true

	case "etcd":
		endpoint = "discovery:///" + grpcClientCfg.Name // Connecting to grpc services by service name
		cli, err := etcdcli.Init(config.Get().Etcd.Addrs, etcdcli.WithDialTimeout(time.Second*2))
		if err != nil {
			panic(err)
		}
		iDiscovery = etcd.New(cli)
		isUseDiscover = true
	case "nacos":
		// example: endpoint = "discovery:///serverName.scheme"
		endpoint = "discovery:///" + grpcClientCfg.Name + ".grpc"
		cli, err := nacoscli.NewNamingClient(
			config.Get().NacosRd.IPAddr,
			config.Get().NacosRd.Port,
			config.Get().NacosRd.NamespaceID)
		if err != nil {
			panic(err)
		}
		iDiscovery = nacos.New(cli)
		isUseDiscover = true

	default:
		endpoint = grpcClientCfg.Host + ":" + strconv.Itoa(grpcClientCfg.Port)
		iDiscovery = nil
		isUseDiscover = false
	}

	if iDiscovery != nil {
		cliOptions = append(cliOptions, grpccli.WithDiscovery(iDiscovery))
	}

	msg := "dialing grpc server"
	if isUseDiscover {
		msg += " with discovery from " + grpcClientCfg.RegistryDiscoveryType
	}
	logger.Info(msg, logger.String("name", grpcClientCfg.Name), logger.String("endpoint", endpoint))

	conn, err := grpccli.Dial(context.Background(), endpoint, cliOptions...)
	if err != nil {
		panic(err)
	}

	return conn
}
