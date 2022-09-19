## grpccli

grpc 客户端，支持服务发现、日志、负载均衡、链路跟踪、指标、重试、熔断。

### 使用示例

```go
func grpcClientExample() pb.UserExampleServiceClient {
	err := config.Init(third_party.Path("../config/conf.yml"))
	if err != nil {
		panic(err)
	}

	var discovery registry.Discovery
	var endpoint = fmt.Sprintf("127.0.0.1:%d", config.Get().Grpc.Port)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	if config.Get().App.EnableRegistryDiscovery {
		endpoint = "discovery:///" + config.Get().App.Name
		discovery = discoveryETCD(config.Get().Etcd.Addrs)
	}
	conn, err := grpccli.DialInsecure(ctx, endpoint,
		grpccli.WithEnableLog(logger.Get()),
		grpccli.WithDiscovery(discovery),
		//grpccli.WithEnableLog(logger.Get()),
		//grpccli.WithDiscovery(discovery),
		//grpccli.WithEnableTrace(),
		//grpccli.WithEnableHystrix("user"),
		//grpccli.WithEnableLoadBalance(),
		//grpccli.WithEnableRetry(),
		//grpccli.WithEnableMetrics(),
	)
	if err != nil {
		panic(err)
	}

	return pb.NewUserExampleServiceClient(conn)
}

func discoveryETCD(endpoints []string) registry.Discovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 10 * time.Second,
		DialOptions: []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	})
	if err != nil {
		panic(err)
	}

	return etcd.New(cli)
}
```
