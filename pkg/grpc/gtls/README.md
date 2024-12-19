## gtls

`gtls` provides grpc secure connectivity by tls, supporting both one-way secure connection and mutual tls connection.

### Example of use

#### One-way secure connection

**grpc server example**

```go
import "github.com/go-dev-frame/sponge/pkg/grpc/gtls"

func main() {
    // one-way connection
    credentials, err := gtls.GetServerTLSCredentials(
        certfile.Path("/one-way/server.crt"),
        certfile.Path("/one-way/server.key"),
    )
    // check err

    server := grpc.NewServer(grpc.Creds(credentials))
}
```

<br>

**grpc client example**

```go
import "github.com/go-dev-frame/sponge/pkg/grpc/gtls"

func main() {
    // one-way connection
    credentials, err := gtls.GetClientTLSCredentials(
        "localhost",
        certfile.Path("/one-way/server.crt"),
	)
    // check err

    conn, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(credentials))
    // check err
}
```

<br>

#### Mutual tls connection

**grpc server example**

```go
import "github.com/go-dev-frame/sponge/pkg/grpc/gtls"

func main() {
    // two-way secure connection
    credentials, err := gtls.GetServerTLSCredentialsByCA(
        certfile.Path("two-way/ca.pem"),
        certfile.Path("two-way/server/server.pem"),
        certfile.Path("two-way/server/server.key"),
    )
    // check err

    server := grpc.NewServer(grpc.Creds(credentials))
}
```

<br>

**grpc client example**

```go
import "github.com/go-dev-frame/sponge/pkg/grpc/gtls"

func main() {
    // two-way secure connection
    credentials, err := gtls.GetClientTLSCredentialsByCA(
        "localhost",
        certfile.Path("two-way/ca.pem"),
        certfile.Path("two-way/client/client.pem"),
        certfile.Path("two-way/client/client.key"),
    )
    // check err

    conn, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(credentials))
    // check err
}
```
