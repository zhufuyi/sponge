## ratelimit

ratelimit 只适用于linux系统。

<br>

### 使用示例

**gin ratelimit 限流中间件示例**

```go
func RateLimit(opts ...RateLimitOption) gin.HandlerFunc {
	o := defaultRatelimitOptions()
	o.apply(opts...)
	limiter := bbr.NewLimiter(
		bbr.WithWindow(o.window),
		bbr.WithBucket(o.bucket),
		bbr.WithCPUThreshold(o.cpuThreshold),
		bbr.WithCPUQuota(o.cpuQuota),
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
