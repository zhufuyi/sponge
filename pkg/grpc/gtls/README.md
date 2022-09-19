## gtls

#### 使用示例

#### grpc server

```go
func main() {
	// 单向认证(服务端认证)
	//credentials, err := gtls.GetServerTLSCredentials(certfile.Path("/one-way/server.crt"), certfile.Path("/one-way/server.key"))

	// 双向认证
	credentials, err := gtls.GetServerTLSCredentialsByCA(
		certfile.Path("two-way/ca.pem"),
		certfile.Path("two-way/server/server.pem"),
		certfile.Path("two-way/server/server.key"),
	)
	if err != nil {
		panic(err)
	}

	// 拦截器
	opts := []grpc.ServerOption{
		grpc.Creds(credentials),
	}

	// 创建grpc server对象，拦截器可以在这里注入
	server := grpc.NewServer(opts...)

	// ......
}
```

<br>

#### grpc client

```go
func main() {
	// 单向认证
	//credentials, err := gtls.GetClientTLSCredentials("localhost", certfile.Path("/one-way/server.crt"))

	// 双向认证
	credentials, err := gtls.GetClientTLSCredentialsByCA(
		"localhost",
		certfile.Path("two-way/ca.pem"),
		certfile.Path("two-way/client/client.pem"),
		certfile.Path("two-way/client/client.key"),
	)
	if err != nil {
		panic(err)
	}

	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(credentials))
	if err != nil {
		panic(err)
	}

	// ......
}
```