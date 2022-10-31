## consulcli

连接etcd服务客户端。

### 使用示例

```go
	endpoints := []string{"192.168.3.37:2379"}
    cli, err := consulcli.Init(endpoints,
        WithConnectTimeout(time.Second*2),
        // WithAuth("", ""),
        // WithAutoSyncInterval(0),
        // WithLog(zap.NewNop()),
	)
```
