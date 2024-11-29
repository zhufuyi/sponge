## keepalive

### Example of use

#### grpc server

```go
func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption
	options = append(options, keepalive.ServerKeepAlive()...)

	return options
}

func main() {
	server := grpc.NewServer(getServerOptions()...)

	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

    // ......
}
```

<br>

#### grpc client

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// use insecure transfer
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// keepalive option
	options = append(options, keepalive.ClientKeepAlive())

	return options
}

func main() {
	conn, err := grpc.NewClient("127.0.0.1:8080", getDialOptions()...)
	if err != nil {
		panic(err)
	}

	// ......
}
```