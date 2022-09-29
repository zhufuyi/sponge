## ratelimiter

gin `path` or `ip` limit.

<br>

### Usage

```go
	r := gin.Default()

	// e.g. (1) default path limit, qps=500, burst=1000
	// r.Use(QPS())

	// e.g. (2) path limit, qps=50, burst=100
	r.Use(QPS(
		WithPath(),
		WithQPS(50),
		WithBurst(100),
	))

	// e.g. (3) ip limit, qps=20, burst=40
	//	r.Use(QPS(
	//		WithIP(),
	//		WithQPS(20),
	//		WithBurst(40),
	//	))
```
