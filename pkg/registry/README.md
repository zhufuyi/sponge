## registry

registry 服务注册，与服务注册[discovery](../discovery)对应，支持etcd、consul、nacos三种方式。

### 使用示例

#### etcd

```go
func setETCDRegistry(etcdEndpoints []string) registry.Registry {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
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

func registryExample() {
    // endpoint组成schema://ip，例如grpc://127.0.0.1:9090, http://127.0.0.1:8080
    endpoints:=[]string{"grpc://127.0.0.1:8080"}
	serviceInstance := registry.NewServiceInstance("user", endpoints,
		//registry.WithID("1"), // 服务id，唯一
		//registry.WithVersion("v0.0.1"), // 服务版本
		//registry.WithMetadata(map[string]string{"foo":"bar"}), // 附加元数据
	)

    iRegistry := setETCDRegistry([]string("192.128.1.6:2349"))

    ctx := context.Background

    // 注册
    ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
    if err := iRegistry.Register(ctx, serviceInstance); err != nil {
        panic(err)
    }

    // 取消注册
    ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)
    if err := iRegistry.Deregister(ctx, serviceInstance); err != nil {
        return err
    }
}
```

<br>

#### consul

<br>

#### nacos


