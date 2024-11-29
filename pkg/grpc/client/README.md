## grpc client

Generic grpc client.

### Example of use

```go
	import "github.com/zhufuyi/sponge/pkg/grpc/client"

	conn, err := client.NewClient(context.Background(), "127.0.0.1:8282",
		//client.WithServiceDiscover(builder),
		//client.WithLoadBalance(),
		//client.WithSecure(credentials),
		//client.WithUnaryInterceptor(unaryInterceptors...),
		//client.WithStreamInterceptor(streamInterceptors...),
	)
```

Examples of practical use https://github.com/zhufuyi/grpc_examples/blob/main/usage/client/main.go
