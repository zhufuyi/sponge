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
			"please change to the correct service name", serverName))
	}

	// 如果没有使用服务发现，用ip和端口直连rpc服务
	endpoint := fmt.Sprintf("%s:%d", grpcClientCfg.Host, cfg.Grpc.Port)

	switch grpcClientCfg.RegistryDiscoveryType {
	// 使用consul发现服务
	case "consul":
		endpoint = "discovery:///" + grpcClientCfg.Name // 通过服务名称连接grpc服务
		cli, err := consulcli.Init(cfg.Consul.Addr, consulcli.WithWaitTime(time.Second*5))
		if err != nil {
			panic(fmt.Sprintf("consulcli.Init error: %v, addr: %s", err, cfg.Consul.Addr))
		}
		iDiscovery := consul.New(cli)
		cliOptions = append(cliOptions, grpccli.WithDiscovery(iDiscovery))
	// 使用etcd发现服务
	case "etcd":
		endpoint = "discovery:///" + grpcClientCfg.Name // 通过服务名称连接grpc服务
		cli, err := etcdcli.Init(cfg.Etcd.Addrs, etcdcli.WithDialTimeout(time.Second*5))
		if err != nil {
			panic(fmt.Sprintf("etcdcli.Init error: %v, addr: %s", err, cfg.Etcd.Addrs))
		}
		iDiscovery := etcd.New(cli)
		cliOptions = append(cliOptions, grpccli.WithDiscovery(iDiscovery))
	// 使用nacos发现服务
	case "nacos":
		// example: endpoint = "discovery:///serverName.scheme"
		endpoint = "discovery:///" + grpcClientCfg.Name + ".grpc"
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

	if cfg.App.EnableTracing {
		cliOptions = append(cliOptions, grpccli.WithEnableTrace())
	}
	if cfg.App.EnableCircuitBreaker {
		cliOptions = append(cliOptions, grpccli.WithEnableCircuitBreaker())
	}
	if cfg.App.EnableMetrics {
		cliOptions = append(cliOptions, grpccli.WithEnableMetrics())
	}

	// 如果需要安全连接，使用grpccli.Dial(ctx, endpoint, cliOptions...)，并且cliOptions设置WithCredentials指定证书路径
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
