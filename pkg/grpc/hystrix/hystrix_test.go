package hystrix

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnaryClientInterceptor(t *testing.T) {
	interceptor := UnaryClientInterceptor("hystrix",
		WithStatsDCollector("localhost:5555", "hystrix", 0.5, 2048))
	assert.NotNil(t, interceptor)
}

func TestStreamClientInterceptor(t *testing.T) {
	interceptor := StreamClientInterceptor("hystrix",
		WithStatsDCollector("localhost:5555", "hystrix", 0.5, 2048))
	assert.NotNil(t, interceptor)
}

func Test_durationToInt(t *testing.T) {
	durationToInt(10*time.Second, time.Second)
}
