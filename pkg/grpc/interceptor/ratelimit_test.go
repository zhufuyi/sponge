package interceptor

import (
	"testing"
	"time"

	"github.com/reugn/equalizer"
	"github.com/stretchr/testify/assert"
)

func TestStreamServerRateLimit(t *testing.T) {
	interceptor := StreamServerRateLimit()
	assert.NotNil(t, interceptor)
}

func TestUnaryServerRateLimit(t *testing.T) {
	interceptor := UnaryServerRateLimit()
	assert.NotNil(t, interceptor)
}

func TestWithRateLimitQPS(t *testing.T) {
	testData := 1000
	opt := WithRateLimitQPS(testData)
	o := new(rateLimitOptions)
	o.apply(opt)
	assert.Less(t, time.Duration(testData), o.refillInterval)

	_ = WithRateLimitQPS(5)
	_ = WithRateLimitQPS(55)
	_ = WithRateLimitQPS(255)
	_ = WithRateLimitQPS(555)
	_ = WithRateLimitQPS(1555)
	_ = WithRateLimitQPS(2555)
	_ = WithRateLimitQPS(5555)
	_ = WithRateLimitQPS(55555)
}

func Test_defaultRateLimitOptions(t *testing.T) {
	o := defaultRateLimitOptions()
	assert.NotNil(t, o)
}

func Test_rateLimitOptions_apply(t *testing.T) {
	testData := 1000
	opt := WithRateLimitQPS(testData)
	o := new(rateLimitOptions)
	o.apply(opt)
	assert.Less(t, time.Duration(testData), o.refillInterval)
}

func Test_myLimiter_Limit(t *testing.T) {
	l := &myLimiter{equalizer.NewTokenBucket(100, 50)}
	actual := l.Limit()
	assert.Equal(t, false, actual)
}
