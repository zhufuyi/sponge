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
// set unary server logging
func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption
	
	options = append(options, grpc_middleware.WithUnaryServerChain(
		// if you don't want to log reply data, you can use interceptor.StreamServerSimpleLog instead of interceptor.UnaryServerLog,
		interceptor.UnaryServerLog(
			logger.Get(),
			interceptor.WithReplaceGRPCLogger(),
			//interceptor.WithMarshalFn(fn), // customised marshal function, default is jsonpb.Marshal
			//interceptor.WithLogIgnoreMethods(fullMethodNames), // ignore methods logging
			//interceptor.WithMaxLen(400), // logging max length, default 300
		),
	))

	return options
}


// you can also set stream server logging
```

**grpc client-side**

```go
// set unary client logging
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	option := grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			interceptor.UnaryClientLog(logger.Get()),
			interceptor.WithReplaceGRPCLogger(),
		),
	)
	options = append(options, option)

	return options
}

// you can also set stream client logging
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

	// rate limiter
	option := grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			interceptor.UnaryRateLimit(
				//interceptor.WithWindow(time.Second*5),
				//interceptor.WithBucket(200),
				//interceptor.WithCPUThreshold(600),
				//interceptor.WithCPUQuota(0), 
			),
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
			interceptor.UnaryClientCircuitBreaker(
				//interceptor.WithValidCode(codes.DeadlineExceeded), // add error code 4 for circuit breaker
				//interceptor.WithUnaryServerDegradeHandler(handler), // add custom degrade handler
			),
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

#### jwt authentication

**grpc client-side**

```go
package main

import (
	"context"
	"github.com/zhufuyi/sponge/pkg/jwt"
	"github.com/zhufuyi/sponge/pkg/grpc/interceptor"
	"github.com/zhufuyi/sponge/pkg/grpc/grpccli"
	userV1 "user/api/user/v1"
)

func main() {
	ctx := context.Background()
	conn, _ := grpccli.Dial(ctx, "127.0.0.1:8282")
	cli := userV1.NewUserClient(conn)

	token := "xxxxxx"
	ctx = interceptor.SetJwtTokenToCtx(ctx, token)

	req := &userV1.GetUserByIDRequest{Id: 100}
	cli.GetByID(ctx, req)
}
```

**grpc server-side**

```go
package main

import (
	"context"
	"net"
	"github.com/zhufuyi/sponge/pkg/jwt"
	"github.com/zhufuyi/sponge/pkg/grpc/interceptor"
	"google.golang.org/grpc"
	userV1 "user/api/user/v1"
)

func standardVerifyFn(claims *jwt.Claims, tokenTail32 string) error {
	// you can check the claims and tokenTail32, and return an error if necessary
	// see example in jwtAuth_test.go line 23

	return nil
}

func customVerifyFn(claims *jwt.CustomClaims, tokenTail32 string) error {
	// you can check the claims and tokenTail32, and return an error if necessary
	// see example in jwtAuth_test.go line 34

	return nil
}

func getUnaryServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// other interceptors ...

	options = append(options, grpc.UnaryInterceptor(
	    // jwt authorization interceptor
	    interceptor.UnaryServerJwtAuth(
			// // choose a verification method as needed
			//interceptor.WithStandardVerify(standardVerifyFn), // standard verify (default), standardVerifyFn is not mandatory, you can set nil if you don't need it
			//interceptor.WithCustomVerify(customVerifyFn), // custom verify

	        // specify the grpc API to ignore token verification(full path)
	        interceptor.WithAuthIgnoreMethods(
	            "/api.user.v1.User/Register",
	            "/api.user.v1.User/Login",
	        ),
	    ),
	))

	return options
}


type user struct {
	userV1.UnimplementedUserServer
}

// Login ...
func (s *user) Login(ctx context.Context, req *userV1.LoginRequest) (*userV1.LoginReply, error) {
	// check user and password success

	uid := 100
	name := "tom"
	token, err := jwt.GenerateToken(uid, name)

	return &userV1.LoginReply{Token: token},nil
}

// GetByID ...
func (s *user) GetByID(ctx context.Context, req *userV1.GetUserByIDRequest) (*userV1.GetUserByIDReply, error) {
	// if token is valid, won't get here, because the interceptor has returned an error message 

	// if you want get jwt claims, you can use the following code
	claims, err := interceptor.GetJwtClaims(ctx)

	return &userV1.GetUserByIDReply{},nil
}

func main()  {
	list, err := net.Listen("tcp", ":8282")
	server := grpc.NewServer(getUnaryServerOptions()...)
	userV1.RegisterUserServer(server, &user{})
	server.Serve(list)
	select{}
}
```

<br>
