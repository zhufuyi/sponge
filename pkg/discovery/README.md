## discovery

discovery 服务服务发现，与服务注册[registry](../registry)对应，支持etcd、consul、nacos三种方式。

### 使用示例

#### etcd

```go
func getETCDDiscovery(etcdEndpoints []string) registry.Discovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: etcdEndpoints,
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

func discoveryExample() {
    iDiscovery := getETCDDiscovery([]string{"192.168.1.6:2379"})

    // 服务发现的endpoint固定格式discovery:///serviceName
	endpoint := "discovery:///" + "user"
    // grpc客户端
	conn, err := grpccli.DialInsecure(ctx, endpoint,
		grpccli.WithUnaryInterceptors(interceptor.UnaryClientLog(logger.Get())),
		grpccli.WithDiscovery(discovery),
	)
	if err != nil {
		panic(err)
	}

	// ......
}
```

<br>

#### consul

<br>

#### nacos


