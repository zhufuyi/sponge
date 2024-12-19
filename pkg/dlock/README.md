## dlock

`dlock` is a distributed lock library based on [**redsync**](https://github.com/go-redsync/redsync) and [**etcd**](https://github.com/etcd-io/etcd). It provides a simple and easy-to-use API for acquiring and releasing locks.

<br>

### Example of use

#### Redis Lock

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/go-dev-frame/sponge/pkg/goredis"
    "github.com/go-dev-frame/sponge/pkg/dlock"
)

func main() {
    redisCli, err := goredis.Init("default:123456@192.168.3.37:6379")  // single redis instance
    // clusterRedisCli, err := goredis.InitCluster(...)                // or init redis cluster
    // redisCli, err := goredis.InitSentinel(...)                      // or init redis sentinel
    if err != nil {
        panic(err)
    }
    defer redisCli.Close()

    locker, err := NewRedisLock(redisCli, "test_lock")
    if err != nil {
        panic(err)
    }
    ctx, _ := context.WithTimeout(context.Background(), time.Second*10)

    // case 1: try to acquire lock, unblock if failed
    {
        ok, err := locker.TryLock(ctx)
        if err != nil {
            panic(err)
        }
        if !ok {
            fmt.Println("failed to acquire lock")
            return
        }
        defer func() {
            if err := locker.Unlock(ctx); err != nil {
                panic(err)
            }
        }()
        // do something here
    }

    // case 2: lock acquired, block until released, timeout, ctx error
    {
        if err := locker.Lock(ctx); err != nil {
            panic(err)
        }
        defer func() {
            if err := locker.Unlock(ctx); err != nil {
                panic(err)
            }
        }()
        // do something here
    }
}
```

<br>

#### Etcd Lock

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/go-dev-frame/sponge/pkg/etcdcli"
    "github.com/go-dev-frame/sponge/pkg/dlock"
)

func main() {
    endpoints := []string{"192.168.3.37:2379"}
    cli, err := etcdcli.Init(endpoints,  etcdcli.WithConnectTimeout(time.Second*5))
    if err!= nil {
        panic(err)
    }
    defer cli.Close()

    locker, err := dlock.NewEtcd(cli, "sponge/dlock", 10)
    if err != nil {
        panic(err)
    }
    ctx, _ := context.WithTimeout(context.Background(), time.Second*10)

    // case 1: try to acquire lock, unblock if failed
    {
        ok, err := locker.TryLock(ctx)
        if err != nil {
            panic(err)
        }
        if !ok {
            fmt.Println("failed to acquire lock")
            return
        }
        defer func() {
            if err := locker.Unlock(ctx); err != nil {
                panic(err)
            }
        }()
        // do something here
    }

    // case 2: lock acquired, block until released, timeout, ctx error
    {
        if err := locker.Lock(ctx); err != nil {
            panic(err)
        }
        defer func() {
            if err := locker.Unlock(ctx); err != nil {
                panic(err)
            }
        }()
        // do something here
    }
}
```
