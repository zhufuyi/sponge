## interceptor

常用grpc客户端和服务端的拦截器。

<br>

### 使用示例

#### jwt

```go
// grpc server

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// token鉴权
	options = append(options, grpc.UnaryInterceptor(
	    interceptor.UnaryServerJwtAuth(
	        // middleware.WithAuthClaimsName("tokenInfo"), // 设置附加到ctx的鉴权信息名称，默认是tokenInfo
	        middleware.WithAuthIgnoreMethods( // 添加忽略token验证的方法
	            "/proto.Account/Register",
	        ),
	    ),
	))

	return options
}

// 生成鉴权信息authorization
func (a *Account) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
    // ......
	token, err := jwt.GenerateToken(uid)
	// handle err
	authorization = middleware.GetAuthorization(token) // 加上前缀
    // ......
}

// 客户端调用方法时必须把鉴权信息通过context传递进来，key名称必须是authorization
func getUser(client pb.AccountClient, req *pb.RegisterReply) error {
	md := metadata.Pairs("authorization", req.Authorization)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: req.Id})
	if err != nil {
		return err
	}

	fmt.Println("get user success", resp)
	return nil
}
```

<br>

#### logging

```go
var logger *zap.Logger

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// 日志设置，默认打印客户端断开连接信息，示例 https://pkg.go.dev/github.com/grpc-ecosystem/go-grpc-middleware/logging/zap
	options = append(options, grpc_middleware.WithUnaryServerChain(
		interceptor.UnaryServerCtxTags(),
		interceptor.UnaryServerZapLogging(
			logger.Get(), // zap
			// middleware.WithLogFields(map[string]interface{}{"serverName": "userExample"}), // 附加打印字段
			middleware.WithLogIgnoreMethods("/proto.userExampleService/GetByID"), // 忽略指定方法打印，可以指定多个
		),
	))

	return options
}
```

<br>

#### recovery

```go
func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	recoveryOption := grpc_middleware.WithUnaryServerChain(
		interceptor.UnaryServerRecovery(),
	)
	options = append(options, recoveryOption)

	return options
}
```

<br>

#### retry

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 禁用tls
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 重试
	option := grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			interceptor.UnaryClientRetry(
                //middleware.WithRetryTimes(5), // 修改默认重试次数，默认3次
                //middleware.WithRetryInterval(100*time.Millisecond), // 修改默认重试时间间隔，默认50毫秒
                //middleware.WithRetryErrCodes(), // 添加触发重试错误码，默认codes.Internal, codes.DeadlineExceeded, codes.Unavailable
			),
		),
	)
	options = append(options, option)

	return options
}
```

<br>

#### 限流

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 禁用tls
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// circuit breaker
	option := grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			interceptor.UnaryRateLimit(),
		),
	)
	options = append(options, option)

	return options
}
```

<br>


#### 熔断器

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 禁用tls
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// circuit breaker
	option := grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			interceptor.UnaryClientCircuitBreaker(),
		),
	)
	options = append(options, option)

	return options
}
```

<br>

#### timeout

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 禁止tls加密
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 超时拦截器
	option := grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			middleware.ContextTimeout(time.Second), // //设置超时
		),
	)
	options = append(options, option)

	return options
}
```

<br>

#### tracing

```go
// 初始化trace
func InitTrace(serviceName string) {
	exporter, err := tracer.NewJaegerAgentExporter("192.168.3.37", "6831")
	if err != nil {
		panic(err)
	}

	resource := tracer.NewResource(
		tracer.WithServiceName(serviceName),
		tracer.WithEnvironment("dev"),
		tracer.WithServiceVersion("demo"),
	)

	tracer.Init(exporter, resource) // 默认采集全部
}

// 在客户端设置链路跟踪
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 禁用tls加密
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// tracing跟踪
	options = append(options, grpc.WithUnaryInterceptor(
		interceptor.UnaryClientTracing(),
	))

	return options
}

// 在服务端设置链路跟踪
func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// 链路跟踪拦截器
	options = append(options, grpc.UnaryInterceptor(
		interceptor.UnaryServerTracing(),
	))

	return options
}

// 如果有需要，可以在程序创建一个span
func SpanDemo(serviceName string, spanName string, ctx context.Context) {
	_, span := otel.Tracer(serviceName).Start(
		ctx, spanName,
		trace.WithAttributes(attribute.String(spanName, time.Now().String())), // 自定义属性
	)
	defer span.End()

	// ......
}
```

<br>

#### metrics

使用示例 [metrics](../metrics/README.md)。

<br>
