## resolve

### Example of use

#### grpc client

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// use insecure transfer
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// load balancing policy, polling https://github.com/grpc/grpc-go/tree/master/examples/features/load_balancing
	// https://github.com/grpc/grpc/blob/master/doc/service_config.md
	options = append(options, grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`))

	return options
}

func main() {
	endpoint := resolve.Register("grpc", "hello.grpc.io", []string{"127.0.0.1:8282", "127.0.0.1:8284", "127.0.0.1:8286"})
	roundRobinConn, err := grpc.Dial(endpoint, getDialOptions()...)
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
