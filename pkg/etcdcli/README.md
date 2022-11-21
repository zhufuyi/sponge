## etcdcli

Connect to the etcd service client.

### Example of use

```go
	endpoints := []string{"192.168.3.37:2379"}
    cli, err := Init(endpoints,
        WithConnectTimeout(time.Second*2),
        // WithAuth("", ""),
        // WithAutoSyncInterval(0),
        // WithLog(zap.NewNop()),
	)
)
```
