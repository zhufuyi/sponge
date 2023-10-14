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

#### common authorization

```go
import "github.com/zhufuyi/sponge/pkg/jwt"

func main() {
    r := gin.Default()

    r.POST("/user/login", Login)
    r.GET("/user/:id", middleware.Auth(), h.GetByID)
	// r.GET("/user/:id", middleware.Auth(middleware.WithVerify(verify)), userFun) // with verify

    r.Run(serverAddr)
}

func verify(claims *jwt.Claims) error {
    if claims.UID != "123" || claims.Role != "admin" {
        return errors.New("verify failed")
    }
    return nil
}

func Login(c *gin.Context) {
	// login success

	// generate token
	token, err := jwt.GenerateToken("123", "admin")
    // handle err
}
```
<br>

#### custom authorization

```go
import "github.com/zhufuyi/sponge/pkg/jwt"

func main() {
    r := gin.Default()

	r.POST("/user/login", Login)
	r.GET("/user/:id", middleware.AuthCustom(verifyCustom), userFun)

    r.Run(serverAddr)
}

func verifyCustom(claims *jwt.CustomClaims) error {
	err := errors.New("verify failed")

	id, exist := claims.Get("id")
	if !exist {
		return err
	}
	foo, exist := claims.Get("foo")
	if !exist {
		return err
	}
	if int(id.(float64)) != fields["id"].(int) || foo.(string) != fields["foo"].(string) {
		return err
	}

	return nil
}

func Login(c *gin.Context) {
    // login success

	// generate token
	fields := jwt.KV{"id": 123, "foo": "bar"}
	token, err := jwt.GenerateCustomToken(fields)
    // handle err
}
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
