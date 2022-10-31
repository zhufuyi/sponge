package initial

import (
	"fmt"
	"strconv"
	"time"

	"github.com/zhufuyi/sponge/internal/config"
	"github.com/zhufuyi/sponge/internal/server"

	"github.com/zhufuyi/sponge/pkg/app"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/consul"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/etcd"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/nacos"
)

// RegisterServers 注册app服务
func RegisterServers() []app.IServer {
	var cfg = config.Get()
	var servers []app.IServer

	// 创建grpc服务
	grpcAddr := ":" + strconv.Itoa(cfg.Grpc.Port)
	grpcRegistry, grpcInstance := registryService("grpc", cfg.App.Host, cfg.Grpc.Port)
	grpcServer := server.NewGRPCServer(grpcAddr,
		server.WithGrpcReadTimeout(time.Duration(cfg.Grpc.ReadTimeout)*time.Second),
		server.WithGrpcWriteTimeout(time.Duration(cfg.Grpc.WriteTimeout)*time.Second),
		server.WithGrpcRegistry(grpcRegistry, grpcInstance),
	)
	servers = append(servers, grpcServer)

	return servers
}

func registryService(scheme string, host string, port int) (registry.Registry, *registry.ServiceInstance) {
	instanceEndpoint := fmt.Sprintf("%s://%s:%d", scheme, host, port)
	cfg := config.Get()

	switch cfg.App.RegistryDiscoveryType {
	// 使用consul注册服务
	case "consul":
		iRegistry, instance, err := consul.NewRegistry(
			cfg.Consul.Addr,
			cfg.App.Name+"_"+scheme+"_"+host,
			cfg.App.Name,
			[]string{instanceEndpoint},
		)
		if err != nil {
			panic(err)
		}
		return iRegistry, instance
	// 使用etcd注册服务
	case "etcd":
		iRegistry, instance, err := etcd.NewRegistry(
			cfg.Etcd.Addrs,
			cfg.App.Name+"_"+scheme+"_"+host,
			cfg.App.Name,
			[]string{instanceEndpoint},
		)
		if err != nil {
			panic(err)
		}
		return iRegistry, instance
	// 使用nacos注册服务
	case "nacos":
		iRegistry, instance, err := nacos.NewRegistry(
			cfg.NacosRd.IPAddr,
			cfg.NacosRd.Port,
			cfg.NacosRd.NamespaceID,
			cfg.App.Name+"_"+scheme+"_"+host,
			cfg.App.Name,
			[]string{instanceEndpoint},
		)
		if err != nil {
			panic(err)
		}
		return iRegistry, instance
	}

	return nil, nil
}
