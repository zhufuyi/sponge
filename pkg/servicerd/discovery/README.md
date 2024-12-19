## discovery

Service discovery, corresponding to the service [registry](../registry), supports etcd, consul and nacos.

### Example of use

```go
    import "github.com/go-dev-frame/sponge/pkg/servicerd/discovery"

    var cliOptions = []grpccli.Option{}
    var endpoint string

	switch grpcClientCfg.RegistryDiscoveryType {
	// discovering services using consul
	case "consul":
		endpoint = "discovery:///" + grpcClientCfg.Name // Connecting to grpc services by service name
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
			panic(fmt.Sprintf("etcdcli.Init error: %v, addr: %s", err, cfg.Etcd.Addrs))
		}
		iDiscovery := etcd.New(cli)
		cliOptions = append(cliOptions, grpccli.WithDiscovery(iDiscovery))
	// discovering services using nacos
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

    serverNameExampleConn, err = grpccli.DialInsecure(context.Background(), endpoint, cliOptions...)
    if err != nil {
        panic(fmt.Sprintf("dial rpc server failed: %v, endpoint: %s", err, endpoint))
    }
```