## loadbalance

### 使用示例

#### grpc client

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 禁止tls加密
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 负载均衡策略，轮询，https://github.com/grpc/grpc-go/tree/master/examples/features/load_balancing 和 https://github.com/grpc/grpc/blob/master/doc/service_config.md
	options = append(options, grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`))

	return options
}

func main() {
	target := loadbalance.Register("grpc", "hello.grpc.io", []string{"127.0.0.1:8080", "127.0.0.1:8081", "127.0.0.1:8082"})
	fmt.Println(target)

	roundRobinConn, err := grpc.Dial(target, getDialOptions()...)
	if err != nil {
		panic(err)
	}
	defer roundRobinConn.Close()

	client := serverNameV1.NewGreeterClient(roundRobinConn)
	for {
		err = sayHello(client)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 2)
	}
}
```