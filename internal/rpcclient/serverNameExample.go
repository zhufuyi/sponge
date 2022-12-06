package rpcclient

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/zhufuyi/sponge/internal/config"

	"github.com/zhufuyi/sponge/pkg/consulcli"
	"github.com/zhufuyi/sponge/pkg/etcdcli"
	"github.com/zhufuyi/sponge/pkg/grpc/grpccli"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/nacoscli"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/consul"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/etcd"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/nacos"

	"google.golang.org/grpc"
)

var (
	serverNameExampleConn *grpc.ClientConn
	serverNameExampleOnce sync.Once
)

// NewServerNameExampleRPCConn instantiate rpc client connection
func NewServerNameExampleRPCConn() {
	cfg := config.Get()

	var cliOptions = []grpccli.Option{
		grpccli.WithEnableLog(logger.Get()),
		//grpccli.WithEnableLoadBalance(),
		//grpccli.WithEnableRetry(),
	}

	serverName := "serverNameExample"
	var grpcClientCfg config.GrpcClient
	for _, cli := range cfg.GrpcClient {
		if strings.EqualFold(cli.Name, serverName) {
			grpcClientCfg = cli
			break
		}
	}
	if grpcClientCfg.Name == "" {
		panic(fmt.Sprintf("not found server name '%v' in yaml config file (field GrpcClient), "+
			"please change to the correct server name", serverName))
	}

	// If service discovery is not used, connect directly to the rpc service using the ip and port
	endpoint := fmt.Sprintf("%s:%d", grpcClientCfg.Host, grpcClientCfg.Port)

	switch grpcClientCfg.RegistryDiscoveryType {
	// discovering services using consul
	case "consul":
		endpoint = "discovery:///" + grpcClientCfg.Name // connecting to grpc services by service name
		cli, err := consulcli.Init(cfg.Consul.Addr, consulcli.WithWaitTime(time.Second*5))
		if err != nil {
			panic(fmt.Sprintf("consulcli.Init error: %v, addr: %s", err, cfg.Consul.Addr))
		}
		iDiscovery := consul.New(cli)
		cliOptions = append(cliOptions, grpccli.WithDiscovery(iDiscovery))
	// discovering services using etcd
	case "etcd":
		endpoint = "discovery:///" + grpcClientCfg.Name // Connecting to grpc services by service name
		cli, err := etcdcli.Init(cfg.Etcd.Addrs, etcdcli.WithDialTimeout(time.Second*5))
		if err != nil {
			panic(fmt.Sprintf("etcdcli.Init error: %v, addr: %v", err, cfg.Etcd.Addrs))
		}
		iDiscovery := etcd.New(cli)
		cliOptions = append(cliOptions, grpccli.WithDiscovery(iDiscovery))
	// discovering services using nacos
	case "nacos":
		// example: endpoint = "discovery:///serverName.scheme"
		endpoint = "discovery:///" + grpcClientCfg.Name + ".grpc" // Connecting to grpc services by service name
		cli, err := nacoscli.NewNamingClient(
			cfg.NacosRd.IPAddr,
			cfg.NacosRd.Port,
			cfg.NacosRd.NamespaceID)
		if err != nil {
			panic(fmt.Sprintf("nacoscli.NewNamingClient error: %v, ipAddr: %s, port: %d",
				err, cfg.NacosRd.IPAddr, cfg.NacosRd.Port))
		}
		iDiscovery := nacos.New(cli)
		cliOptions = append(cliOptions, grpccli.WithDiscovery(iDiscovery))
	}

	if cfg.App.EnableTrace {
		cliOptions = append(cliOptions, grpccli.WithEnableTrace())
	}
	if cfg.App.EnableCircuitBreaker {
		cliOptions = append(cliOptions, grpccli.WithEnableCircuitBreaker())
	}
	if cfg.App.EnableMetrics {
		cliOptions = append(cliOptions, grpccli.WithEnableMetrics())
	}

	// If a secure connection is required, use grpccli.Dial(ctx, endpoint, cliOptions...) and
	// cliOptions sets WithCredentials to specify the certificate path
	var err error
	serverNameExampleConn, err = grpccli.DialInsecure(context.Background(), endpoint, cliOptions...)
	if err != nil {
		panic(fmt.Sprintf("dial rpc server failed: %v, endpoint: %s", err, endpoint))
	}
}

// GetServerNameExampleRPCConn get client conn
func GetServerNameExampleRPCConn() *grpc.ClientConn {
	if serverNameExampleConn == nil {
		serverNameExampleOnce.Do(func() {
			NewServerNameExampleRPCConn()
		})
	}

	return serverNameExampleConn
}

// CloseServerNameExampleRPCConn Close tears down the ClientConn and all underlying connections.
func CloseServerNameExampleRPCConn() error {
	if serverNameExampleConn == nil {
		return nil
	}

	return serverNameExampleConn.Close()
}
