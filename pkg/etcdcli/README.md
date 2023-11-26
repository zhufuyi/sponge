## etcdcli

Connect to the etcd service client.

### Example of use

```go
    import "github.com/zhufuyi/sponge/pkg/etcdcli"

    endpoints := []string{"192.168.3.37:2379"}
    // Way 1: setting parameters
    cli, err := etcdcli.Init(
        endpoints,
        etcdcli.WithConnectTimeout(time.Second*2),
        // etcdcli.WithAutoSyncInterval(0),
        // etcdcli.WithLog(zap.NewNop()),
        // etcdcli.WithAuth("", ""),
    )

    // Way 2: Setting up clientv3.Config
    cli, err = etcdcli.Init(nil, etcdcli.WithConfig(&clientv3.Config{
        Endpoints:   endpoints,
        DialTimeout: time.Second * 2,
        //Username:    "",
        //Password:    "",
    }))
```
