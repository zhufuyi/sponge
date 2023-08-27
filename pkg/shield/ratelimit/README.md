## ratelimit

Adaptive rate limit, only available for linux systems.

<br>

### Example of use

#### gin ratelimit middleware

```go
import (
rl "github.com/zhufuyi/sponge/pkg/shield/ratelimit"
)

func RateLimit(opts ...RateLimitOption) gin.HandlerFunc {
	o := defaultRatelimitOptions()
	o.apply(opts...)
	limiter := rl.NewLimiter(
		rl.WithWindow(o.window),
		rl.WithBucket(o.bucket),
		rl.WithCPUThreshold(o.cpuThreshold),
		rl.WithCPUQuota(o.cpuQuota),
	)

	return func(c *gin.Context) {
		done, err := limiter.Allow()
		if err != nil {
			response.Output(c, http.StatusTooManyRequests, err.Error())
			c.Abort()
			return
		}

		c.Next()

		done(rl.DoneInfo{Err: c.Request.Context().Err()})
	}
}
```


<br>

#### grpc ratelimit interceptor

```go
import (
	rl "github.com/zhufuyi/sponge/pkg/shield/ratelimit"
)


func UnaryServerRateLimit(opts ...RatelimitOption) grpc.UnaryServerInterceptor {
	o := defaultRatelimitOptions()
	o.apply(opts...)
	limiter := rl.NewLimiter(
		rl.WithWindow(o.window),
		rl.WithBucket(o.bucket),
		rl.WithCPUThreshold(o.cpuThreshold),
		rl.WithCPUQuota(o.cpuQuota),
	)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		done, err := limiter.Allow()
		if err != nil {
			return nil, errcode.StatusLimitExceed.ToRPCErr(err.Error())
		}

		reply, err := handler(ctx, req)
		done(rl.DoneInfo{Err: err})
		return reply, err
	}
}
```