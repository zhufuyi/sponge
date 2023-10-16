## etcdcli

Connect to the etcd service client.

### Example of use

```go
    import "github.com/zhufuyi/sponge/pkg/etcdcli"

    endpoints := []string{"192.168.3.37:2379"}
    cli, err := etcdcli.Init(
        endpoints,
        WithConnectTimeout(time.Second*2),
        // WithAutoSyncInterval(0),
        // WithLog(zap.NewNop()),
        // WithAuth("", ""),
    )
```
