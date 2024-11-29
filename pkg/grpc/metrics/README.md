## metrics

The grpc's server-side and client-side metrics can continue to be captured using prometheus.

### Example of use

#### grpc server

```go
import "github.com/zhufuyi/sponge/pkg/grpc/metrics"

func UnaryServerLabels(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// set up prometheus custom labels
	tag := grpc_ctxtags.NewTags().
		Set(serverNameLabelKey, serverNameLabelValue).
		Set(envLabelKey, envLabelValue)
	newCtx := grpc_ctxtags.SetInContext(ctx, tag)

	return handler(newCtx, req)
}

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// metrics interceptor
	option := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		//UnaryServerLabels,                  // tag
		metrics.UnaryServerMetrics(
			// metrics.WithCounterMetrics(customizedCounterMetric) // adding custom metrics
		),
	))
	options = append(options, option)

	option = grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		metrics.StreamServerMetrics(), // metrics interceptor for streaming rpc
	))
	options = append(options, option)

	return options
}

func main() {
	rand.Seed(time.Now().UnixNano())

	addr := ":8282"
	fmt.Println("start rpc server", addr)

	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(getServerOptions()...)
	serverNameV1.RegisterGreeterServer(server, &GreeterServer{})

	// start metrics server, collect grpc metrics by default, turn on, go metrics
	metrics.ServerHTTPService(":8283", server)
	fmt.Println("start metrics server", ":8283")

	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
```

<br>

#### grpc client

```go
import "github.com/zhufuyi/sponge/pkg/grpc/metrics"

func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// use insecure transfer
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// Metrics
	options = append(options, grpc.WithUnaryInterceptor(metrics.UnaryClientMetrics()))
	options = append(options, grpc.WithStreamInterceptor(metrics.StreamClientMetrics()))
	return options
}

func main() {
	conn, err := grpc.NewClient("127.0.0.1:8282", getDialOptions()...)

	metrics.ClientHTTPService(":8284")
	fmt.Println("start metrics server", ":8284")

	client := serverNameV1.NewGreeterClient(conn)
	i := 0
	for {
		i++
		time.Sleep(time.Millisecond * 500) // qps is 2
		err = sayHello(client, i)
		if err != nil {
			fmt.Println(err)
		}
	}
}
```