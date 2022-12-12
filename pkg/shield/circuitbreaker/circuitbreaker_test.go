package circuitbreaker_test

import (
	"testing"

	"github.com/zhufuyi/sponge/pkg/shield/circuitbreaker"
)

// This is an example of using a circuit breaker Do() when return nil.
func TestCircuitBreaker(t *testing.T) {
	b := circuitbreaker.NewBreaker()
	for i := 0; i < 1000; i++ {
		b.MarkSuccess()
	}
	for i := 0; i < 100; i++ {
		b.MarkFailed()
	}

	err := b.Allow()
	t.Log(err)
	// Output: err=<nil>
}
