## middleware

Gin middleware plugin.

<br>

## Example of use

### logging middleware

You can set the maximum length for printing, add a request id field, ignore print path, customize [zap](go.uber.org/zap) log.

```go
    import "github.com/zhufuyi/sponge/pkg/gin/middleware"

    r := gin.Default()

    // default
    r.Use(middleware.Logging()) // simplified logging using middleware.SimpleLog()

    // --- or ---

    // custom
    r.Use(middleware.Logging(    // simplified logging using middleware.SimpleLog(WithRequestIDFromHeader())
        middleware.WithMaxLen(400),
        WithRequestIDFromHeader(),
        //WithRequestIDFromContext(),
        //middleware.WithLog(log), // custom zap log
        //middleware.WithIgnoreRoutes("/hello"),
    ))
```

<br>

### Allow cross-domain requests middleware

```go
    import "github.com/zhufuyi/sponge/pkg/gin/middleware"

    r := gin.Default()
    r.Use(middleware.Cors())
```

<br>

### rate limiter middleware

Adaptive flow limitation based on hardware resources.

```go
    import "github.com/zhufuyi/sponge/pkg/gin/middleware"

    r := gin.Default()

    // default
    r.Use(middleware.RateLimit())

    // --- or ---

    // custom
    r.Use(middleware.RateLimit(
        WithWindow(time.Second*10),
        WithBucket(100),
        WithCPUThreshold(100),
        WithCPUQuota(0.5),
    ))
```

<br>

### Circuit Breaker middleware

```go
    import "github.com/zhufuyi/sponge/pkg/gin/middleware"

    r := gin.Default()
    r.Use(middleware.CircuitBreaker(
        //middleware.WithValidCode(http.StatusRequestTimeout), // add error code 408 for circuit breaker
        //middleware.WithDegradeHandler(handler), // add custom degrade handler
    ))
```

<br>

### jwt authorization middleware

#### standard authorization

```go
import "github.com/zhufuyi/sponge/pkg/jwt"
import "github.com/zhufuyi/sponge/pkg/gin/middleware"

func main() {
    r := gin.Default()

    r.POST("/user/login", Login)
    r.GET("/user/:id", middleware.Auth(), h.GetByID) // do not get claims
    // r.GET("/user/:id", middleware.Auth(middleware.WithVerify(adminVerify)), h.GetByID) // get claims and check

    r.Run(serverAddr)
}

func adminVerify(claims *jwt.Claims, tokenTail10 string, c *gin.Context) error {
    if claims.Name != "admin" {
        return errors.New("verify failed")
    }

    // token := getToken(claims.UID) // from cache or database
    // if tokenTail10 != token[len(token)-10:] { return err }

    return nil
}

func Login(c *gin.Context) {
    // generate token
    token, err := jwt.GenerateToken("123", "admin")
    // save token
}
```

<br>

#### custom authorization

```go
import "github.com/zhufuyi/sponge/pkg/jwt"
import "github.com/zhufuyi/sponge/pkg/gin/middleware"

func main() {
    r := gin.Default()

    r.POST("/user/login", Login)
    r.GET("/user/:id", middleware.AuthCustom(verify), h.GetByID) // get claims and check

    r.Run(serverAddr)
}

// custom verify example
func verify(claims *jwt.CustomClaims, tokenTail10 string, c *gin.Context) error {
    err := errors.New("verify failed")

    token, fields := getToken(id) // from cache or database
    // if tokenTail10 != token[len(token)-10:] { return err }

    id, exist := claims.GetUint64("id")
    if !exist || id != fields["id"].(uint64) { return err }

    name, exist := claims.GetString("name")
    if !exist || name != fields["name"].(string) { return err }

    age, exist := claims.GetInt("age")
    if !exist || age != fields["age"].(int) { return err }

    return nil
}

func Login(c *gin.Context) {
    // generate token
    fields := jwt.KV{"id": uint64(123), "name": "tom", "age": 10}
    token, err := jwt.GenerateCustomToken(fields)
    // save token and fields
}
```

<br>

### tracing middleware

```go
import "github.com/zhufuyi/sponge/pkg/tracer"
import "github.com/zhufuyi/sponge/pkg/gin/middleware"

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
    import "github.com/zhufuyi/sponge/pkg/gin/middleware/metrics"

    r := gin.Default()

    r.Use(metrics.Metrics(r,
        //metrics.WithMetricsPath("/demo/metrics"), // default is /metrics
        metrics.WithIgnoreStatusCodes(http.StatusNotFound), // ignore status codes
        //metrics.WithIgnoreRequestMethods(http.MethodHead),  // ignore request methods
        //metrics.WithIgnoreRequestPaths("/ping", "/health"), // ignore request paths
    ))
```

<br>

### Request id

```go
    import "github.com/zhufuyi/sponge/pkg/gin/middleware"

    // Default request id
    r := gin.Default()
    r.Use(middleware.RequestID())

    // --- or ---

    // Customized request id key
    //r.User(middleware.RequestID(
    //    middleware.WithContextRequestIDKey("your ctx request id key"), // default is request_id
    //    middleware.WithHeaderRequestIDKey("your header request id key"), // default is X-Request-Id
    //))
    // If you change the ContextRequestIDKey, you have to set the same key name if you want to print the request id in the mysql logs as well.
    // example: 
    // db, err := mysql.Init(dsn,
        // mysql.WithLogRequestIDKey("your ctx request id key"),  // print request_id
        // ...
    // )
```

<br>

### Timeout

```go
    import "github.com/zhufuyi/sponge/pkg/gin/middleware"

    r := gin.Default()

    // way1: global set timeout
    r.Use(middleware.Timeout(time.Second*5))

    // --- or ---

    // way2: router set timeout
    r.GET("/userExample/:id", middleware.Timeout(time.Second*3), h.GetByID)

    // Note: If timeout is set both globally and in the router, the minimum timeout prevails
```
