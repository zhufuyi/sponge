## goredis

在 [go-redis]github.com/go-redis/redis 基础上封装的库。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/goredis

<br>

## 使用示例

```go
	redisCli, err := goredis.Init(config.Get().RedisURL, goredis.WithEnableTrace())
	if err != nil {
		panic("goredis.Init error: " + err.Error())
	}
```

<br>

官方文档 https://redis.uptrace.dev/guide/
