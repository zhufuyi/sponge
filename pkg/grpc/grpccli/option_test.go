package grpccli

import (
	"go.uber.org/zap"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/grpc/interceptor"
	"github.com/zhufuyi/sponge/pkg/registry"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestWithCredentials(t *testing.T) {
	testData := insecure.NewCredentials()
	opt := WithCredentials(testData)
	o := new(options)
	o.apply(opt)
	assert.Equal(t, testData, o.credentials)
}

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

func TestWithEnableHystrix(t *testing.T) {
	testData := "hystrix"
	opt := WithEnableHystrix(testData)
	o := new(options)
	o.apply(opt)
	assert.Equal(t, testData, o.hystrixName)
}

func TestWithEnableLoadBalance(t *testing.T) {
	opt := WithEnableLoadBalance()
	o := new(options)
	o.apply(opt)
	assert.Equal(t, true, o.enableLoadBalance)
}

func TestWithEnableLog(t *testing.T) {
	testData := zap.NewNop()
	opt := WithEnableLog(testData)
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
	assert.Equal(t, testData, o.timeout)
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
