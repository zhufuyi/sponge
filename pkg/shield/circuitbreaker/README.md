## circuitbreaker

circuitbreaker 熔断器。

<br>

### 使用示例

**gin circuitbreaker 熔断器中间件示例**

```go
func CircuitBreaker(opts ...CircuitBreakerOption) gin.HandlerFunc {
	o := defaultCircuitBreakerOptions()
	o.apply(opts...)

	return func(c *gin.Context) {
		breaker := o.group.Get(c.FullPath()).(circuitbreaker.CircuitBreaker)
		if err := breaker.Allow(); err != nil {
			// NOTE: when client reject request locally,
			// continue to add counter let the drop ratio higher.
			breaker.MarkFailed()
			response.Output(c, http.StatusServiceUnavailable, err.Error())
			c.Abort()
			return
		}

		c.Next()

		code := c.Writer.Status()
		// NOTE: need to check internal and service unavailable error
		if code == http.StatusInternalServerError || code == http.StatusServiceUnavailable || code == http.StatusGatewayTimeout {
			breaker.MarkFailed()
		} else {
			breaker.MarkSuccess()
		}
	}
}
```
