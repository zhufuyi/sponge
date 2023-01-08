## middleware

gin middleware plugin.

<br>

## Example of use

### logging middleware

You can set the maximum length for printing, add a request id field, ignore print path, customize [zap](go.uber.org/zap) log

```go
    r := gin.Default()
	
    r.Use(middleware.Logging())

    r.Use(middleware.Logging(
        middleware.WithMaxLen(400),
		//WithRequestIDFromHeader(),
		WithRequestIDFromContext(),
        //middleware.WithIgnoreRoutes("/hello"),
    ))

    log, _ := logger.Init(logger.WithFormat("json"))
    r.Use(middlewareLogging(
        middleware.WithLog(log),
    ))
```

<br>

### Allow cross-domain requests middleware

```go
    r := gin.Default()
    r.Use(middleware.Cors())
```

<br>

### rate limiter middleware

Adaptive flow limitation based on hardware resources.

```go
	r := gin.Default()

    // e.g. (1) use default
    // r.Use(RateLimit())
    
    // e.g. (2) custom parameters
    r.Use(RateLimit(
    WithWindow(time.Second*10),
    WithBucket(100),
    WithCPUThreshold(100),
    WithCPUQuota(0.5),
    ))
```

<br>

### Circuit Breaker middleware

```go
    r := gin.Default()
    r.Use(CircuitBreaker())
```
<br>

### jwt authorization middleware

```go
    r := gin.Default()
    r.GET("/user/:id", middleware.JWT(), userFun)
```
<br>

### tracing middleware

```go
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

	tracer.Init(exporter, resource) // collect all by default
}

func NewRouter(
    r := gin.Default()
    r.Use(middleware.Tracing("your-service-name"))

    // ......
)

// if necessary, you can create a span in the program
func SpanDemo(serviceName string, spanName string, ctx context.Context) {
	_, span := otel.Tracer(serviceName).Start(
		ctx, spanName,
		trace.WithAttributes(attribute.String(spanName, time.Now().String())),
	)
	defer span.End()

	// ......
}
```

<br>

### Metrics middleware

```go
	r := gin.Default()

	r.Use(metrics.Metrics(r,
		//metrics.WithMetricsPath("/demo/metrics"), // default is /metrics
		metrics.WithIgnoreStatusCodes(http.StatusNotFound), // ignore status codes
		//metrics.WithIgnoreRequestMethods(http.MethodHead),  // ignore request methods
		//metrics.WithIgnoreRequestPaths("/ping", "/health"), // ignore request paths
	))
```

### Request id

```go
	r := gin.Default()
    r.Use(RequestID())
```
