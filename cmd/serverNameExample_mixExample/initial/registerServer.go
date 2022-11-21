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

// RegisterServers register for the app service
func RegisterServers() []app.IServer {
	var cfg = config.Get()
	var servers []app.IServer

	// creating http service
	httpAddr := ":" + strconv.Itoa(cfg.HTTP.Port)
	httpRegistry, httpInstance := registryService("http", cfg.App.Host, cfg.HTTP.Port)
	httpServer := server.NewHTTPServer(httpAddr,
		server.WithHTTPReadTimeout(time.Second*time.Duration(cfg.HTTP.ReadTimeout)),
		server.WithHTTPWriteTimeout(time.Second*time.Duration(cfg.HTTP.WriteTimeout)),
		server.WithHTTPRegistry(httpRegistry, httpInstance),
		server.WithHTTPIsProd(cfg.App.Env == "prod"),
	)
	servers = append(servers, httpServer)

	// creating grpc service
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
	// registering service with consul
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
	// registering service with etcd
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
	// registering service with nacos
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
