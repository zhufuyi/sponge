## consulcli

Connect to the consul service client.

### Example of use

```go
    import "github.com/zhufuyi/sponge/pkg/consulcli"

    addr := "192.168.3.37:8500"

    // Way 1: setting parameters
    cli, err := consulcli.Init(addr,
        consulcli.WithWaitTime(time.Second*5),
        // consulcli.WithDatacenter(""),
    )

    // Way 2: setting up api.Config
    cli, err = Init("", consulcli.WithConfig(&api.Config{
        Address:    addr,
        Scheme:     "http",
        WaitTime:   time.Second * 5,
        Datacenter: "",
    }))
```
