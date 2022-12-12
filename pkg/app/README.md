## app

Start and stop services gracefully, using [errgroup](golang.org/x/sync/errgroup) to ensure that multiple services are started properly at the same time.

<br>

### Example of use

```go
func main() {
    initApp()
    servers := registerServers()
    closes := registerCloses(servers)

    a := app.New(servers, closes)
    a.Run()
}

func initApp() {
    // get configuration

    // initializing log

    // initializing database

    // ......
}

func registerServers() []app.IServer {
    var servers []app.IServer

    // creating http service
    servers = append(servers, server.NewHTTPServer(

    ))

    // creating grpc service
    servers = append(servers, server.NewGRPCServer(

    ))

    // ......

    return servers
}

func registerCloses(servers []app.IServer) []app.Close {
    var closes []app.Close

    // close server
    for _, server := range servers {
        closes = append(closes, server.Stop)
    }

    // close other resource
    closes = append(closes, func() error {

    })

    // ......

    return closes
}
```
