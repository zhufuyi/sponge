## goredis

`goredis` library wrapped in [go-redis](github.com/go-redis/redis).

<br>

## Example of use

```go
	redisCli, err := goredis.Init(config.Get().RedisURL, goredis.WithEnableTrace())
	if err != nil {
		panic("goredis.Init error: " + err.Error())
	}
```

<br>

Official Documents https://redis.uptrace.dev/guide/
