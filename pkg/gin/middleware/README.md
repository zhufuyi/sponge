## render

gin中间件插件。

<br>

## 使用示例

### 日志中间件

可以设置打印最大长度、添加请求id字段、忽略打印path、自定义[zap](go.uber.org/zap) log

```go
    r := gin.Default()

    // 默认打印日志
    r.Use(middleware.Logging())

    // 自定义打印日志
    r.Use(middleware.Logging(
        middleware.WithMaxLen(400), // 打印body最大长度，超过则忽略
		//WithRequestIDFromHeader(), // 支持自定义requestID名称
		WithRequestIDFromContext(), // 支持自定义requestID名称
        //middleware.WithIgnoreRoutes("/hello"), // 忽略/hello
    ))

    // 自定义zap log
    log, _ := logger.Init(logger.WithFormat("json"))
    r.Use(middlewareLogging(
        middleware.WithLog(log),
    ))
```

<br>

### 允许跨域请求

```go
    r := gin.Default()
    r.Use(middleware.Cors())
```

<br>

### qps限流

#### path维度的qps限流

```go
    r := gin.Default()

    // path, 默认qps=500, burst=1000
    r.Use(ratelimiter.QPS())

    // path, 自定义qps=50, burst=100
    r.Use(ratelimiter.QPS(
        ratelimiter.WithQPS(50),
        ratelimiter.WithBurst(100),
    ))
```

#### ip维度的qps限流

```go
    // ip, 自定义qps=40, burst=80
    r.Use(ratelimiter.QPS(
        ratelimiter.WithIP(),
        ratelimiter.WithQPS(40),
        ratelimiter.WithBurst(80),
    ))
```

<br>

### jwt鉴权

```go
    r := gin.Default()
    r.GET("/user/:id", middleware.JWT(), userFun) // 需要鉴权
```
<br>

### 链路跟踪

```go
// 初始化trace
func InitTrace(serviceName string) {
	exporter, err := tracer.NewJaegerAgentExporter("192.168.3.37", "6831")
	if err != nil {
		panic(err)
	}

	resource := tracer.NewResource(
		tracer.WithServiceName(serviceName),
		tracer.WithEnvironment("dev"),
		tracer.WithServiceVersion("demo"),
	)

	tracer.Init(exporter, resource) // 默认采集全部
}

func NewRouter(
    r := gin.Default()
    r.Use(middleware.Tracing("your-service-name"))

    // ......
)

// 如果有需要，可以在程序创建一个span
func SpanDemo(serviceName string, spanName string, ctx context.Context) {
	_, span := otel.Tracer(serviceName).Start(
		ctx, spanName,
		trace.WithAttributes(attribute.String(spanName, time.Now().String())), // 自定义属性
	)
	defer span.End()

	// ......
}
```

<br>

### 监控指标

```go
	r := gin.Default()

	r.Use(metrics.Metrics(r,
		//metrics.WithMetricsPath("/demo/metrics"), // default is /metrics
		metrics.WithIgnoreStatusCodes(http.StatusNotFound), // ignore status codes
		//metrics.WithIgnoreRequestMethods(http.MethodHead),  // ignore request methods
		//metrics.WithIgnoreRequestPaths("/ping", "/health"), // ignore request paths
	))
```