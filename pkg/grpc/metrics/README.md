## metrics

grpc的server和client端指标，可以使用prometheus继续采集这些指标。

### 使用示例

#### grpc server

```go
func UnaryServerLabels(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 设置prometheus公共标签
	tag := grpc_ctxtags.NewTags().
		Set(serverNameLabelKey, serverNameLabelValue).
		Set(envLabelKey, envLabelValue)
	newCtx := grpc_ctxtags.SetInContext(ctx, tag)

	return handler(newCtx, req)
}

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// metrics拦截器
	option := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		//UnaryServerLabels,                  // 标签
		metrics.UnaryServerMetrics(
			// metrics.WithCounterMetrics(customizedCounterMetric) // 添加自定义指标
		), // 一元rpc的metrics拦截器
	))
	options = append(options, option)

	option = grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		metrics.StreamServerMetrics(), // 流式rpc的metrics拦截器
	))
	options = append(options, option)

	return options
}

func main() {
	rand.Seed(time.Now().UnixNano())

	addr := ":8080"
	fmt.Println("start rpc server", addr)

	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(getServerOptions()...)
    serverNameV1.RegisterGreeterServer(server, &GreeterServer{})

	// 启动metrics服务器，默认采集grpc指标，开启、go指标
	metrics.ServerHTTPService(":9092", server)
	fmt.Println("start metrics server", ":9092")

	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
```

<br>

#### grpc client

```go
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
	conn, err := grpc.Dial("127.0.0.1:8080", getDialOptions()...)

	metrics.ClientHTTPService(":9094")
	fmt.Println("start metrics server", ":9094")

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