## interceptor

Commonly used grpc client-side and server-side interceptors.

<br>

### Example of use

```go
import "github.com/zhufuyi/sponge/pkg/grpc/interceptor"
```

#### logging

**grpc server-side**

```go
var logger *zap.Logger

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption
	
	options = append(options, grpc_middleware.WithUnaryServerChain(
		interceptor.UnaryClientLog(
			logger.Get(), // zap
			// middleware.WithLogFields(map[string]interface{}{"serverName": "userExample"}), // additional print fields
			middleware.WithLogIgnoreMethods("/proto.userExampleService/GetByID"), // ignore the specified method print, you can specify more than one
		),
	))

	return options
}
```

**grpc client-side**

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	option := grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			interceptor.UnaryClientLog(logger.Get()),
		),
	)
	options = append(options, option)

	return options
}
```

<br>

#### recovery

**grpc server-side**

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

**grpc client-side**

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	option := grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			interceptor.UnaryClientRecovery(),
		),
	)
	options = append(options, option)

	return options
}
```

<br>

#### retry

**grpc client-side**

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// use insecure transfer
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// retry
	option := grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			interceptor.UnaryClientRetry(
				//middleware.WithRetryTimes(5), // modify the default number of retries to 3 by default
				//middleware.WithRetryInterval(100*time.Millisecond), // modify the default retry interval, default 50 milliseconds
				//middleware.WithRetryErrCodes(), // add trigger retry error code, default is codes.Internal, codes.DeadlineExceeded, codes.Unavailable
			),
		),
	)
	options = append(options, option)

	return options
}
```

<br>

#### rate limiter

**grpc server-side**

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// use insecure transfer
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

#### Circuit Breaker

**grpc server-side**

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// use insecure transfer
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

**grpc client-side**

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// use insecure transfer
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// timeout
	option := grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			middleware.ContextTimeout(time.Second), // set timeout
		),
	)
	options = append(options, option)

	return options
}
```

<br>

#### tracing

**grpc server-side**

```go
// initialize tracing
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

	tracer.Init(exporter, resource) // collect all by default
}

// set up trace on the client side
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// use insecure transfer
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// use tracing
	options = append(options, grpc.WithUnaryInterceptor(
		interceptor.UnaryClientTracing(),
	))

	return options
}

// set up trace on the server side
func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// use tracing
	options = append(options, grpc.UnaryInterceptor(
		interceptor.UnaryServerTracing(),
	))

	return options
}

// if necessary, you can create a span in the program
func SpanDemo(serviceName string, spanName string, ctx context.Context) {
	_, span := otel.Tracer(serviceName).Start(
		ctx, spanName,
		trace.WithAttributes(attribute.String(spanName, time.Now().String())), // customised attributes
	)
	defer span.End()

	// ......
}
```

<br>

#### metrics

example [metrics](../metrics/README.md).

<br>

#### Request id

**grpc server-side**

```go
func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	recoveryOption := grpc_middleware.WithUnaryServerChain(
		interceptor.UnaryServerRequestID(),
	)
	options = append(options, recoveryOption)

	return options
}
```

<br>

**grpc client-side**

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// use insecure transfer
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	option := grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			interceptor.UnaryClientRequestID(),
		),
	)
	options = append(options, option)

	return options
}
```

<br>

#### jwt

**grpc server-side**

```go
func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// token authorization
	options = append(options, grpc.UnaryInterceptor(
	    interceptor.UnaryServerJwtAuth(
	        // middleware.WithAuthClaimsName("tokenInfo"), // set the name of the forensic information attached to the ctx, the default is tokenInfo
	        middleware.WithAuthIgnoreMethods( // add a way to ignore token validation
	            "/proto.Account/Register",
	        ),
	    ),
	))

	return options
}

// generate forensic information authorization
func (a *Account) Register(ctx context.Context, req *serverNameV1.RegisterRequest) (*serverNameV1.RegisterReply, error) {
	// ......
	token, err := jwt.GenerateToken(uid)
	// handle err
	authorization = middleware.GetAuthorization(token)
	// ......
}

// the client must pass in the authentication information via the context when calling the method, and the key name must be authorization
func getUser(client serverNameV1.AccountClient, req *serverNameV1.RegisterReply) error {
	md := metadata.Pairs("authorization", req.Authorization)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := client.GetUser(ctx, &serverNameV1.GetUserRequest{Id: req.Id})
	if err != nil {
		return err
	}

	fmt.Println("get user success", resp)
	return nil
}
```

<br>
