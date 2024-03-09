package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zhufuyi/sponge/pkg/container/group"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/shield/circuitbreaker"
)

// ErrNotAllowed error not allowed.
var ErrNotAllowed = circuitbreaker.ErrNotAllowed

// CircuitBreakerOption set the circuit breaker circuitBreakerOptions.
type CircuitBreakerOption func(*circuitBreakerOptions)

type circuitBreakerOptions struct {
	group *group.Group
	// http code for circuit breaker, default already includes 500 and 503
	validCodes map[int]struct{}
}

func defaultCircuitBreakerOptions() *circuitBreakerOptions {
	return &circuitBreakerOptions{
		group: group.NewGroup(func() interface{} {
			return circuitbreaker.NewBreaker()
		}),
		validCodes: map[int]struct{}{
			http.StatusInternalServerError: {},
			http.StatusServiceUnavailable:  {},
		},
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
		if g != nil {
			o.group = g
		}
	}
}

// WithValidCode http code to mark failed
func WithValidCode(code ...int) CircuitBreakerOption {
	return func(o *circuitBreakerOptions) {
		for _, c := range code {
			o.validCodes[c] = struct{}{}
		}
	}
}

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
		// NOTE: need to check internal and service unavailable error
		_, isHit := o.validCodes[code]
		if isHit {
			breaker.MarkFailed()
		} else {
			breaker.MarkSuccess()
		}
	}
}
