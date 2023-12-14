## grpc server

Generic grpc server.

### Example of use

```go
	import "github.com/zhufuyi/sponge/pkg/grpc/server"

	port := 8282
	fn := func(s *grpc.Server) {
		pb.RegisterGreeterServer(s, &greeterServer{})
	}
	
	server.Run(port, []server.RegisterFn{fn},
		//server.WithSecure(credentials),
		//server.WithUnaryInterceptor(unaryInterceptors...),
		//server.WithStreamInterceptor(streamInterceptors...),
		//server.WithServiceRegister(func() {}),
	)

	select{}
```

Examples of practical use https://github.com/zhufuyi/grpc_examples/blob/main/wrapGrpc/server/main.go
