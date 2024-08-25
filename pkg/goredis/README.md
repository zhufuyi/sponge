## goredis

`goredis` library wrapped in [go-redis](https://github.com/redis/go-redis).

<br>

### Example of use

#### Single Redis

```go
	// way 1: redis version 6.0 or above
	redisCli, err := goredis.Init("default:123456@127.0.0.1:6379") // can set parameters such as timeout, tls, tracing, such as goredis.Withxxx()
	if err != nil {
		panic("goredis.Init error: " + err.Error())
	}

	// way 2: redis version 5.0 or below
	redisCli, err := goredis.InitSingle("127.0.0.1:6379", "123456", 0) // can set parameters such as timeout, tls, tracing, such as goredis.Withxxx()
```

<br>

#### Sentinel

```go
	addrs := []string{"127.0.0.1:6380", "127.0.0.1:6381", "127.0.0.1:6382","127.0.0.1:26380", "127.0.0.1:26381", "127.0.0.1:26382"}
	rdbCli, err : := goredis.InitSentinel("mymaster", addrs, "", "123456") // can set parameters such as timeout, tls, tracing, such as goredis.Withxxx()
```

<br>

#### Cluster

```go
	addrs := []string{"127.0.0.1:6380", "127.0.0.1:6381", "127.0.0.1:6382","127.0.0.1:6383", "127.0.0.1:6384", "127.0.0.1:6385"}
	clusterRdb, err : := goredis.InitCluster(addrs, "", "123456") // can set parameters such as timeout, tls, tracing, such as goredis.Withxxx()
```

<br>

Official Documents https://redis.uptrace.dev/zh/guide/go-redis.html
