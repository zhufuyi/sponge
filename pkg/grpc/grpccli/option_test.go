package grpccli

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-dev-frame/sponge/pkg/grpc/interceptor"
	"github.com/go-dev-frame/sponge/pkg/servicerd/registry"
)

func TestWithDialOptions(t *testing.T) {
	testData := grpc.WithTransportCredentials(insecure.NewCredentials())
	opt := WithDialOptions(testData)
	o := new(options)
	o.apply(opt)
	assert.Equal(t, testData, o.dialOptions[0])
}

func TestWithDiscovery(t *testing.T) {
	testData := new(registry.Discovery)
	opt := WithDiscovery(*testData)
	o := new(options)
	o.apply(opt)
	assert.NotEqual(t, testData, o.discovery)
}

func TestWithEnableCircuitBreaker(t *testing.T) {
	opt := WithEnableCircuitBreaker()
	o := new(options)
	o.apply(opt)
	assert.Equal(t, true, o.enableCircuitBreaker)
}

func TestWithEnableLoadBalance(t *testing.T) {
	opt := WithEnableLoadBalance()
	o := new(options)
	o.apply(opt)
	assert.Equal(t, true, o.enableLoadBalance)
}

func TestWithEnableRequestID(t *testing.T) {
	opt := WithEnableRequestID()
	o := new(options)
	o.apply(opt)
	assert.Equal(t, true, o.enableRequestID)
}

func TestWithEnableLog(t *testing.T) {
	opt := WithEnableLog(nil)
	testData := zap.NewNop()
	opt = WithEnableLog(testData)
	o := new(options)
	o.apply(opt)
	assert.Equal(t, testData, o.log)
}

func TestWithEnableMetrics(t *testing.T) {
	opt := WithEnableMetrics()
	o := new(options)
	o.apply(opt)
	assert.Equal(t, true, o.enableMetrics)
}

func TestWithEnableRetry(t *testing.T) {
	opt := WithEnableRetry()
	o := new(options)
	o.apply(opt)
	assert.Equal(t, true, o.enableRetry)
}

func TestWithEnableTrace(t *testing.T) {
	opt := WithEnableTrace()
	o := new(options)
	o.apply(opt)
	assert.Equal(t, true, o.enableTrace)
}

func TestWithStreamInterceptors(t *testing.T) {
	testData := interceptor.StreamClientRetry()
	opt := WithStreamInterceptors(testData)
	o := new(options)
	o.apply(opt)
	assert.LessOrEqual(t, 1, len(o.streamInterceptors))
}

func TestWithTimeout(t *testing.T) {
	testData := time.Second
	opt := WithTimeout(testData)
	o := new(options)
	o.apply(opt)
	assert.Equal(t, testData, o.requestTimeout)
}

func TestWithDiscoveryInsecure(t *testing.T) {
	var testData bool
	opt := WithDiscoveryInsecure(testData)
	o := new(options)
	o.apply(opt)
	assert.Equal(t, testData, o.discoveryInsecure)
}

func TestWithUnaryInterceptors(t *testing.T) {
	testData := interceptor.UnaryClientRetry()
	opt := WithUnaryInterceptors(testData)
	o := new(options)
	o.apply(opt)
	assert.LessOrEqual(t, 1, len(o.unaryInterceptors))
}

func Test_defaultOptions(t *testing.T) {
	o := defaultOptions()
	assert.NotNil(t, o)
}

func Test_options_apply(t *testing.T) {
	opt := WithEnableRetry()
	o := new(options)
	o.apply(opt)
	assert.Equal(t, true, o.enableRetry)
}

func Test_options_isSecure(t *testing.T) {
	o := new(options)
	secure := o.isSecure()
	assert.Equal(t, false, secure)
	o.secureType = secureOneWay
	secure = o.isSecure()
	assert.Equal(t, true, secure)
}

func TestWithSecure(t *testing.T) {
	o := new(options)
	opt := WithSecure("foo", "", "", "", "")
	o.apply(opt)
	assert.Equal(t, "foo", o.secureType)

	opt = WithSecure(secureOneWay, "", "", "", "")
	o.apply(opt)
	assert.Equal(t, secureOneWay, o.secureType)

	opt = WithSecure(secureTwoWay, "", "", "", "")
	o.apply(opt)
	assert.Equal(t, secureTwoWay, o.secureType)
}
