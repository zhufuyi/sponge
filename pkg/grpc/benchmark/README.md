## benchmark

Compression testing of rpc methods and generation of reported results.

### Example of use

```go
import "github.com/go-dev-frame/sponge/pkg/grpc/benchmark"

func benchmarkExample() error {
	host := "127.0.0.1:8282"
	protoFile := "api/serverNameExample/v1/userExample.proto"
	// if third-party dependencies are missing during the press test, copy them to the project's third_party directory (not including the import path)
	importPaths := []string{"third_party"}
	message := &serverNameV1.GetUserExampleByIDRequest{
		ID: 2,
	}

	b, err := benchmark.New(host, protoFile, "GetByID", message, 1000, importPaths...)
	if err != nil {
		return err
	}
	return b.Run()
}
```

Once the crush is complete, copy the output html file path to your browser to view the detailed crush report.
