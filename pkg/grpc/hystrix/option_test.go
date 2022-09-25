package hystrix

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithErrorPercentThreshold(t *testing.T) {
	testData := 50
	opt := WithErrorPercentThreshold(testData)
	o := new(options)
	o.apply(opt)
	assert.Equal(t, testData, o.errorPercentThreshold)
}

func TestWithFallbackFunc(t *testing.T) {
	testData := func(ctx context.Context, err error) error {
		t.Log("this is fall back")
		return nil
	}
	opt := WithFallbackFunc(testData)
	o := new(options)
	o.apply(opt)
	assert.Equal(t, nil, o.fallbackFunc(context.Background(), nil))
}

func TestWithMaxConcurrentRequests(t *testing.T) {
	testData := 1000
	opt := WithMaxConcurrentRequests(testData)
	o := new(options)
	o.apply(opt)
	assert.Equal(t, testData, o.maxConcurrentRequests)
}

func TestWithPrometheus(t *testing.T) {
	opt := WithPrometheus()
	o := new(options)
	o.apply(opt)
}

func TestWithRequestVolumeThreshold(t *testing.T) {
	testData := 1000
	opt := WithRequestVolumeThreshold(testData)
	o := new(options)
	o.apply(opt)
	assert.Equal(t, testData, o.requestVolumeThreshold)
}

func TestWithSleepWindow(t *testing.T) {
	testData := time.Second * 10
	opt := WithSleepWindow(testData)
	o := new(options)
	o.apply(opt)
	assert.Equal(t, testData, o.sleepWindow)
}

func TestWithStatsDCollector(t *testing.T) {
	opt := WithStatsDCollector("localhost:5555", "hystrix", 0.5, 2048)
	o := new(options)
	o.apply(opt)
	assert.Equal(t, "hystrix", o.statsD.Prefix)
}

func TestWithTimeout(t *testing.T) {
	testData := time.Second * 10
	opt := WithTimeout(testData)
	o := new(options)
	o.apply(opt)
	assert.Equal(t, testData, o.timeout)
}

func Test_defaultOptions(t *testing.T) {
	o := defaultOptions()
	assert.NotNil(t, o)
}

func Test_options_apply(t *testing.T) {
	testData := time.Second * 10
	opt := WithTimeout(testData)
	o := new(options)
	o.apply(opt)
	assert.Equal(t, testData, o.timeout)
}
