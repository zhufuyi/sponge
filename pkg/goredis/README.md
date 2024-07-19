## goredis

`goredis` library wrapped in [go-redis](github.com/go-redis/redis).

<br>

### Example of use

#### Single Redis

```go
	// way 1: redis version 6.0 or above
	redisCli, err := goredis.Init("default:123456@127.0.0.1:6379", goredis.WithEnableTrace())
	if err != nil {
		panic("goredis.Init error: " + err.Error())
	}

	// way 2: redis version 5.0 or below
	redisCli := goredis.Init2("127.0.0.1:6379", "123456", 0, goredis.WithEnableTrace())
```

<br>

#### Sentinel

```go
	addrs := []string{"127.0.0.1:6380", "127.0.0.1:6381", "127.0.0.1:6382"}
	rdb := goredis.InitSentinel("master", addrs, "default", "123456", goredis.WithEnableTrace())
```

<br>

#### Cluster

```go
	addrs := []string{"127.0.0.1:6380", "127.0.0.1:6381", "127.0.0.1:6382"}
	clusterRdb := goredis.InitCluster(addrs, "default", "123456", goredis.WithEnableTrace())
```

<br>

Official Documents https://redis.uptrace.dev/zh/guide/go-redis.html
