## keepalive

### 使用示例

#### grpc server

```go
func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption
	// 默认是每15秒向客户端发送一次ping，修改为间隔20秒发送一次ping
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
	conn, err := grpc.Dial("127.0.0.1:8080", getDialOptions()...)
	if err != nil {
		panic(err)
	}

	// ......
}
```