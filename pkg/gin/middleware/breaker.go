package middleware

import (
	"net/http"

	"github.com/zhufuyi/sponge/pkg/container/group"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/shield/circuitbreaker"

	"github.com/gin-gonic/gin"
)

// ErrNotAllowed error not allowed.
var ErrNotAllowed = circuitbreaker.ErrNotAllowed

// CircuitBreakerOption set the circuit breaker circuitBreakerOptions.
type CircuitBreakerOption func(*circuitBreakerOptions)

type circuitBreakerOptions struct {
	group *group.Group
}

func defaultCircuitBreakerOptions() *circuitBreakerOptions {
	return &circuitBreakerOptions{
		group: group.NewGroup(func() interface{} {
			return circuitbreaker.NewBreaker()
		}),
	}
}

func (o *circuitBreakerOptions) apply(opts ...CircuitBreakerOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithGroup with circuit breaker group.
// NOTE: implements generics circuitbreaker.CircuitBreaker
func WithGroup(g *group.Group) CircuitBreakerOption {
	return func(o *circuitBreakerOptions) {
		o.group = g
	}
}

// CircuitBreaker a circuit breaker middleware
func CircuitBreaker(opts ...CircuitBreakerOption) gin.HandlerFunc {
	o := defaultCircuitBreakerOptions()
	o.apply(opts...)

	return func(c *gin.Context) {
		breaker := o.group.Get(c.FullPath()).(circuitbreaker.CircuitBreaker)
		if err := breaker.Allow(); err != nil {
			// NOTE: when client reject request locally,
			// continue add counter let the drop ratio higher.
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
