## app

优雅的启动和停止服务，使用[errgroup](golang.org/x/sync/errgroup)保证多个服务同时正常启动。

<br>

### 安装

> go get -u github.com/zhufuyi/pkg/app

<br>

### 使用示例

```go
func main() {
	inits := registerInits()
	servers := registerServers()
	closes := registerCloses(servers)

	s := app.New(inits, servers, closes)
	s.Run()
}

func registerInits() []app.Init {
    // 读取配置文件

    var inits []app.Init

	// 初始化日志
	inits = append(inits, func() {

	})

	// 初始化数据库
	inits = append(inits, func() {

	})

    // ......

	return inits
}

func registerServers() []app.IServer {
	var servers []app.IServer

	// 创建http服务
	servers = append(servers, server.NewHTTPServer(

	))

	// 创建grpc服务
	servers = append(servers, server.NewGRPCServer(

	))

    // ......

	return servers
}

func registerCloses(servers []app.IServer) []app.Close {
	var closes []app.Close

	// 关闭服务
	for _, server := range servers {
		closes = append(closes, server.Stop)
	}

	// 关闭数据库连接
	closes = append(closes, func() error {

	})

	// ......

	return closes
}
```