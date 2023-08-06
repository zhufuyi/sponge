## circuitbreaker

Circuit Breaker for web middleware and rpc interceptor.

<br>

### Example of use

**gin circuit breaker middleware**

```go
// CircuitBreaker a circuit breaker middleware
func CircuitBreaker(opts ...CircuitBreakerOption) gin.HandlerFunc {
	o := defaultCircuitBreakerOptions()
	o.apply(opts...)

	return func(c *gin.Context) {
		breaker := o.group.Get(c.FullPath()).(circuitbreaker.CircuitBreaker)
		if err := breaker.Allow(); err != nil {
			// NOTE: when client reject request locally, keep adding counter let the drop ratio higher.
			breaker.MarkFailed()
			response.Output(c, http.StatusServiceUnavailable, err.Error())
			c.Abort()
			return
		}

		c.Next()

		code := c.Writer.Status()
		// NOTE: need to check internal and service unavailable error, e.g. http.StatusInternalServerError
		_, isHit := o.validCodes[code]
		if isHit {
			breaker.MarkFailed()
		} else {
			breaker.MarkSuccess()
		}
	}
}
```

<br>

**rpc server circuit breaker interceptor** 

```go
// UnaryServerCircuitBreaker server-side unary circuit breaker interceptor
func UnaryServerCircuitBreaker(opts ...CircuitBreakerOption) grpc.UnaryServerInterceptor {
	o := defaultCircuitBreakerOptions()
	o.apply(opts...)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		breaker := o.group.Get(info.FullMethod).(circuitbreaker.CircuitBreaker)
		if err := breaker.Allow(); err != nil {
			// NOTE: when client reject request locally, keep adding let the drop ratio higher.
			breaker.MarkFailed()
			return nil, errcode.StatusServiceUnavailable.ToRPCErr(err.Error())
		}

		reply, err := handler(ctx, req)
		if err != nil {
			// NOTE: need to check internal and service unavailable error
			s, ok := status.FromError(err)
			_, isHit := o.validCodes[s.Code()]
			if ok && isHit {
				breaker.MarkFailed()
			} else {
				breaker.MarkSuccess()
			}
		}

		return reply, err
	}
}
```