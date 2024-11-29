package rpcclient

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/zhufuyi/sponge/internal/config"
	"github.com/zhufuyi/sponge/pkg/grpc/grpccli"
	"github.com/zhufuyi/sponge/pkg/logger"
)

var (
	serverNameExampleConn *grpc.ClientConn
	serverNameExampleOnce sync.Once
)

// NewServerNameExampleRPCConn instantiate rpc client connection
func NewServerNameExampleRPCConn() {
	cfg := config.Get()

	serverName := "serverNameExample"
	var grpcClientCfg config.GrpcClient
	for _, cli := range cfg.GrpcClient {
		if strings.EqualFold(cli.Name, serverName) {
			grpcClientCfg = cli
			break
		}
	}
	if grpcClientCfg.Name == "" {
		panic(fmt.Sprintf("not found grpc service name '%v' in configuration file(yaml), "+
			"please add gprc service configuration in the configuration file(yaml) under the field grpcClient.", serverName))
	}

	var cliOptions = []grpccli.Option{
		grpccli.WithEnableRequestID(),
		grpccli.WithEnableLog(logger.Get()),
	}

	// if service discovery is not used, connect directly to the rpc service using the ip and port
	endpoint := fmt.Sprintf("%s:%d", grpcClientCfg.Host, grpcClientCfg.Port)
	isUseDiscover := false

	// using service discovery
	//discoverOption, discoveryEndpoint := discoverService(cfg, grpcClientCfg)
	//if discoverOption != nil {
	//	isUseDiscover = true
	//	endpoint = discoveryEndpoint
	//	cliOptions = append(cliOptions, discoverOption)
	//	cliOptions = append(cliOptions, grpccli.WithEnableLoadBalance()) // load balance
	//}

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

	if cfg.App.EnableTrace {
		cliOptions = append(cliOptions, grpccli.WithEnableTrace())
	}
	if cfg.App.EnableCircuitBreaker {
		cliOptions = append(cliOptions, grpccli.WithEnableCircuitBreaker())
	}
	if cfg.App.EnableMetrics {
		cliOptions = append(cliOptions, grpccli.WithEnableMetrics())
	}
	if grpcClientCfg.Timeout > 0 {
		cliOptions = append(cliOptions, grpccli.WithTimeout(time.Second*time.Duration(grpcClientCfg.Timeout)))
	}

	msg := "dial grpc server"
	if isUseDiscover {
		msg += " with service discovery from " + grpcClientCfg.RegistryDiscoveryType
	}
	logger.Info(msg, logger.String("name", serverName), logger.String("endpoint", endpoint))

	var err error
	serverNameExampleConn, err = grpccli.NewClient(endpoint, cliOptions...)
	if err != nil {
		panic(fmt.Sprintf("grpccli.NewClient error: %v, name: %s, endpoint: %s", err, serverName, endpoint))
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

// discovery service with consul or etcd or nacos, select one of them to use
//func discoverService(cfg *config.Config, grpcClientCfg config.GrpcClient) (grpccli.Option, string) {
//	var (
//		endpoint      string
//		grpcCliOption grpccli.Option
//	)
//
//	switch grpcClientCfg.RegistryDiscoveryType {
//	case "consul":
//		endpoint = "discovery:///" + grpcClientCfg.Name // format: discovery:///serverName
//		cli, err := consulcli.Init(cfg.Consul.Addr, consulcli.WithWaitTime(time.Second*5))
//		if err != nil {
//			panic(fmt.Sprintf("consulcli.Init error: %v, addr: %s", err, cfg.Consul.Addr))
//		}
//		iDiscovery := consul.New(cli)
//		grpcCliOption = grpccli.WithDiscovery(iDiscovery)
//
//	case "etcd":
//		endpoint = "discovery:///" + grpcClientCfg.Name // format: discovery:///serverName
//		cli, err := etcdcli.Init(cfg.Etcd.Addrs, etcdcli.WithDialTimeout(time.Second*5))
//		if err != nil {
//			panic(fmt.Sprintf("etcdcli.Init error: %v, addr: %v", err, cfg.Etcd.Addrs))
//		}
//		iDiscovery := etcd.New(cli)
//		grpcCliOption = grpccli.WithDiscovery(iDiscovery)
//
//	case "nacos":
//		endpoint = "discovery:///" + grpcClientCfg.Name + ".grpc" // format: discovery:///serverName.scheme
//		cli, err := nacoscli.NewNamingClient(
//			cfg.NacosRd.IPAddr,
//			cfg.NacosRd.Port,
//			cfg.NacosRd.NamespaceID)
//		if err != nil {
//			panic(fmt.Sprintf("nacoscli.NewNamingClient error: %v, ipAddr: %s, port: %d",
//				err, cfg.NacosRd.IPAddr, cfg.NacosRd.Port))
//		}
//		iDiscovery := nacos.New(cli)
//		grpcCliOption = grpccli.WithDiscovery(iDiscovery)
//	}
//
//	return grpcCliOption, endpoint
//}
