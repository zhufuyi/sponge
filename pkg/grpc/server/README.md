## grpc server

Generic grpc server.

### Example of use

```go
	import "github.com/go-dev-frame/sponge/pkg/grpc/server"

	port := 8282
	registerFn := func(s *grpc.Server) {
		pb.RegisterGreeterServer(s, &greeterServer{})
	}
	
	server.Run(port, registerFn,
		//server.WithSecure(credentials),
		//server.WithUnaryInterceptor(unaryInterceptors...),
		//server.WithStreamInterceptor(streamInterceptors...),
		//server.WithServiceRegister(func() {}),
		//server.WithStatConnections(metrics.WithConnectionsLogger(logger.Get()), metrics.WithConnectionsGauge()),  // show connections or set prometheus metrics
	)

	select{}
```

Examples of practical use https://github.com/zhufuyi/grpc_examples/blob/main/usage/server/main.go
