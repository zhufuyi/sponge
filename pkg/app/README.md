## app

Elegantly start and stop services, using [errgroup](golang.org/x/sync/errgroup) to ensure that multiple services are started properly at the same time.

<br>

### Example of use

```go
func main() {
	inits := registerInits()
	servers := registerServers()
	closes := registerCloses(servers)

	s := app.New(inits, servers, closes)
	s.Run()
}

func registerInits() []app.Init {
    // get configuration

    var inits []app.Init

	// initializing log
	inits = append(inits, func() {

	})

	// initializing database
	inits = append(inits, func() {

	})

    // ......

	return inits
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
