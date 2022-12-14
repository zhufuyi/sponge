## registry

Service registry, corresponding to service [discovery](../discovery) corresponds to and supports etcd, consul and nacos.

### Example of use

```go
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

// ------------------------------------------------------------------------------------------

    iRegistry, serviceInstance := registryService("http", "127.0.0.1", 8080)
    
    // register service
    ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
    if err := iRegistry.Register(ctx, serviceInstance); err != nil {
        panic(err)
    }
    
    // deregister service
    ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)
    if err := iRegistry.Deregister(ctx, serviceInstance); err != nil {
        return err
    }
```

