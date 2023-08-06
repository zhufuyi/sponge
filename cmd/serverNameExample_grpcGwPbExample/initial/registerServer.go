package initial

import (
	"fmt"
	"strconv"
	"time"

	"github.com/zhufuyi/sponge/internal/config"
	"github.com/zhufuyi/sponge/internal/server"

	"github.com/zhufuyi/sponge/pkg/app"
	"github.com/zhufuyi/sponge/pkg/logger"
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
	httpServer := server.NewHTTPServer_pbExample(httpAddr,
		server.WithHTTPReadTimeout(time.Second*time.Duration(cfg.HTTP.ReadTimeout)),
		server.WithHTTPWriteTimeout(time.Second*time.Duration(cfg.HTTP.WriteTimeout)),
		server.WithHTTPRegistry(httpRegistry, httpInstance),
		server.WithHTTPIsProd(cfg.App.Env == "prod"),
	)
	servers = append(servers, httpServer)

	return servers
}

func registryService(scheme string, host string, port int) (registry.Registry, *registry.ServiceInstance) {
	var (
		instanceEndpoint = fmt.Sprintf("%s://%s:%d", scheme, host, port)
		cfg              = config.Get()

		iRegistry registry.Registry
		instance  *registry.ServiceInstance
		err       error

		id       = cfg.App.Name + "_" + scheme + "_" + host
		logField logger.Field
	)

	switch cfg.App.RegistryDiscoveryType {
	// registering service with consul
	case "consul":
		iRegistry, instance, err = consul.NewRegistry(
			cfg.Consul.Addr,
			id,
			cfg.App.Name,
			[]string{instanceEndpoint},
		)
		if err != nil {
			panic(err)
		}
		logField = logger.Any("consulAddress", cfg.Consul.Addr)

	// registering service with etcd
	case "etcd":
		iRegistry, instance, err = etcd.NewRegistry(
			cfg.Etcd.Addrs,
			id,
			cfg.App.Name,
			[]string{instanceEndpoint},
		)
		if err != nil {
			panic(err)
		}
		logField = logger.Any("etcdAddress", cfg.Etcd.Addrs)

	// registering service with nacos
	case "nacos":
		iRegistry, instance, err = nacos.NewRegistry(
			cfg.NacosRd.IPAddr,
			cfg.NacosRd.Port,
			cfg.NacosRd.NamespaceID,
			id,
			cfg.App.Name,
			[]string{instanceEndpoint},
		)
		if err != nil {
			panic(err)
		}
		logField = logger.String("nacosAddress", fmt.Sprintf("%v:%d", cfg.NacosRd.IPAddr, cfg.NacosRd.Port))
	}

	if instance != nil {
		msg := fmt.Sprintf("register service address to %s", cfg.App.RegistryDiscoveryType)
		logger.Info(msg, logField, logger.String("id", id), logger.String("name", cfg.App.Name), logger.String("endpoint", instanceEndpoint))
		return iRegistry, instance
	}

	return nil, nil
}
