## gtls

#### Example of use

#### grpc server

```go
func main() {
	// one-way authentication (server-side authentication)
	//credentials, err := gtls.GetServerTLSCredentials(certfile.Path("/one-way/server.crt"), certfile.Path("/one-way/server.key"))

	// two-way authentication
	credentials, err := gtls.GetServerTLSCredentialsByCA(
		certfile.Path("two-way/ca.pem"),
		certfile.Path("two-way/server/server.pem"),
		certfile.Path("two-way/server/server.key"),
	)
	if err != nil {
		panic(err)
	}

	// interceptor
	opts := []grpc.ServerOption{
		grpc.Creds(credentials),
	}

	server := grpc.NewServer(opts...)

	// ......
}
```

<br>

#### grpc client

```go
func main() {
	// one-way authentication
	//credentials, err := gtls.GetClientTLSCredentials("localhost", certfile.Path("/one-way/server.crt"))

	// two-way authentication
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