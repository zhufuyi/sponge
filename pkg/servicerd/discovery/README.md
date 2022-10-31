## discovery

discovery 服务发现，与服务注册[registry](../registry)对应，支持etcd、consul、nacos三种方式。

### 使用示例

```go
    var cliOptions = []grpccli.Option{}
    var endpoint string

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

    serverNameExampleConn, err = grpccli.DialInsecure(context.Background(), endpoint, cliOptions...)
    if err != nil {
        panic(fmt.Sprintf("dial rpc server failed: %v, endpoint: %s", err, endpoint))
    }
```