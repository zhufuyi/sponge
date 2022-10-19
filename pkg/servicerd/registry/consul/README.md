## consul

### 使用示例

```go
func registryExample() {
    addr := "127.0.0.1:8500"
    instanceName := "serverName"
    instanceEndpoints := []string{"grpc://127.0.0.1:8282"}
    iRegistry, serviceInstance, err := NewRegistry(addr, instanceName, instanceEndpoints)
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
