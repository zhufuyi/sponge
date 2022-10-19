## etcd

### 使用示例

**example 1**

```go
func registryExample() {
    etcdEndpoints := []string{"127.0.0.1:2379"}
    instanceName := "serverName"
    instanceEndpoints := []string{"grpc://127.0.0.1:8282"}
    iRegistry, serviceInstance, err := NewRegistry(etcdEndpoints, instanceName, instanceEndpoints)
    if err != nil {
        panic(err)
    }

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

**example 2**

```go
// 连接etcd服务器，并实例化Registry接口
func newETCDRegistry(etcdEndpoints []string) registry.Registry {
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
    // endpoint组成schema://ip，例如grpc://127.0.0.1:8282, http://127.0.0.1:8080
    endpoints:=[]string{"grpc://127.0.0.1:8282"}
	serviceInstance := registry.NewServiceInstance("user", endpoints,
		//registry.WithID("1"), // 服务id，唯一编号
		//registry.WithVersion("v0.0.1"), // 服务版本
		//registry.WithMetadata(map[string]string{"foo":"bar"}), // 附加元数据
	)

    iRegistry := newETCDRegistry([]string("192.128.1.6:2349"))
	
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