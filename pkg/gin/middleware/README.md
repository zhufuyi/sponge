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

    // e.g. (1) use default
    // r.Use(middleware.RateLimit())
    
    // e.g. (2) custom parameters
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
    r.Use(middleware.CircuitBreaker())
```

<br>

### jwt authorization middleware

#### common authorization

```go
import "github.com/zhufuyi/sponge/pkg/jwt"
import "github.com/zhufuyi/sponge/pkg/gin/middleware"

func main() {
    r := gin.Default()

    r.POST("/user/login", Login)
    r.GET("/user/:id", middleware.Auth(), h.GetByID) // no verify field
    // r.GET("/user/:id", middleware.Auth(middleware.WithVerify(adminVerify)), h.GetByID) // with verify field

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
    r.GET("/user/:id", middleware.AuthCustom(verify), h.GetByID)

    r.Run(serverAddr)
}

func verify(claims *jwt.CustomClaims, tokenTail10 string, c *gin.Context) error {
    err := errors.New("verify failed")

    // token, fields := getToken(id) // from cache or database
    // if tokenTail10 != token[len(token)-10:] { return err }
	
    id, exist := claims.Get("id")
    if !exist {
        return err
    }
    foo, exist := claims.Get("foo")
    if !exist {
        return err
    }
    if int(id.(float64)) != fields["id"].(int) ||
        foo.(string) != fields["foo"].(string) {
        return err
    }

    return nil
}

func Login(c *gin.Context) {
    // generate token
    fields := jwt.KV{"id": 123, "foo": "bar"}
    token, err := jwt.GenerateCustomToken(fields)
    // save token end fields
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

### Request id

```go
	import "github.com/zhufuyi/sponge/pkg/gin/middleware"

    r := gin.Default()
    r.Use(middleware.RequestID())

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
