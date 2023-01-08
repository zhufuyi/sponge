package interceptor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestStreamClientRetry(t *testing.T) {
	interceptor := StreamClientRetry()
	assert.NotNil(t, interceptor)
}

func TestUnaryClientRetry(t *testing.T) {
	interceptor := UnaryClientRetry()
	assert.NotNil(t, interceptor)
}

func TestWithRetryErrCodes(t *testing.T) {
	testData := codes.Canceled
	opt := WithRetryErrCodes(testData)
	o := new(retryOptions)
	o.apply(opt)
	assert.Contains(t, o.errCodes, testData)
}

func TestWithRetryInterval(t *testing.T) {
	testData := time.Second
	opt := WithRetryInterval(testData)
	o := new(retryOptions)
	o.apply(opt)
	assert.Equal(t, testData, o.interval)

	testData = time.Microsecond
	opt = WithRetryInterval(testData)
	o = new(retryOptions)
	o.apply(opt)
	assert.Equal(t, true, o.interval == time.Millisecond)

	testData = time.Minute
	opt = WithRetryInterval(testData)
	o = new(retryOptions)
	o.apply(opt)
	assert.Equal(t, true, o.interval == 10*time.Second)
}

func TestWithRetryTimes(t *testing.T) {
	testData := uint(5)
	opt := WithRetryTimes(testData)
	o := new(retryOptions)
	o.apply(opt)
	assert.Equal(t, testData, o.times)

	testData = uint(20)
	opt = WithRetryTimes(testData)
	o = new(retryOptions)
	o.apply(opt)
	assert.NotEqual(t, testData, o.times)
}

func Test_defaultRetryOptions(t *testing.T) {
	o := defaultRetryOptions()
	assert.NotNil(t, o)
}

func Test_retryOptions_apply(t *testing.T) {
	testData := uint(5)
	opt := WithRetryTimes(testData)
	o := new(retryOptions)
	o.apply(opt)
	assert.Equal(t, testData, o.times)
}
