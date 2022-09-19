package hystrix

import (
	"context"
	"time"

	"github.com/gunerhuseyin/goprometheus"
	hystrixmiddleware "github.com/gunerhuseyin/goprometheus/middleware/hystrix"

	"github.com/afex/hystrix-go/plugins"
)

const (
	defaultHystrixTimeout         = 30 * time.Second
	defaultMaxConcurrentRequests  = 100
	defaultErrorPercentThreshold  = 25
	defaultSleepWindow            = 10
	defaultRequestVolumeThreshold = 10

	maxUint = ^uint(0)
	maxInt  = int(maxUint >> 1)
)

func defaultOptions() *options {
	return &options{
		fallbackFunc:           nil,
		timeout:                defaultHystrixTimeout,
		maxConcurrentRequests:  defaultMaxConcurrentRequests,
		errorPercentThreshold:  defaultErrorPercentThreshold,
		sleepWindow:            defaultSleepWindow,
		requestVolumeThreshold: defaultRequestVolumeThreshold,
	}
}

// options is the hystrix client implementation
type options struct {
	timeout                time.Duration
	maxConcurrentRequests  int
	requestVolumeThreshold int
	sleepWindow            time.Duration
	errorPercentThreshold  int
	fallbackFunc           func(ctx context.Context, err error) error
	statsD                 *plugins.StatsdCollectorConfig
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// Option represents the hystrix client options
type Option func(*options)

// WithTimeout sets hystrix timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *options) {
		c.timeout = timeout
	}
}

// WithMaxConcurrentRequests sets hystrix max concurrent requests
func WithMaxConcurrentRequests(maxConcurrentRequests int) Option {
	return func(c *options) {
		c.maxConcurrentRequests = maxConcurrentRequests
	}
}

// WithRequestVolumeThreshold sets hystrix request volume threshold
func WithRequestVolumeThreshold(requestVolumeThreshold int) Option {
	return func(c *options) {
		c.requestVolumeThreshold = requestVolumeThreshold
	}
}

// WithSleepWindow sets hystrix sleep window
func WithSleepWindow(sleepWindow time.Duration) Option {
	return func(c *options) {
		c.sleepWindow = sleepWindow
	}
}

// WithErrorPercentThreshold sets hystrix error percent threshold
func WithErrorPercentThreshold(errorPercentThreshold int) Option {
	return func(c *options) {
		c.errorPercentThreshold = errorPercentThreshold
	}
}

// WithFallbackFunc sets the fallback function
func WithFallbackFunc(fn func(ctx context.Context, err error) error) Option {
	return func(c *options) {
		c.fallbackFunc = fn
	}
}

// WithPrometheus sets the hystrix metrics
func WithPrometheus() Option {
	return func(c *options) {
		gpm := goprometheus.New()
		// 采集hystrix指标
		gpHystrixConfig := &hystrixmiddleware.Config{
			Prefix: "hystrix_circuit_breaker_",
		}
		gpHystrix := hystrixmiddleware.New(gpm, gpHystrixConfig)
		gpm.UseHystrix(gpHystrix)

		gpm.Run() // 添加go prometheus路由/metrics
	}
}

// WithStatsDCollector exports hystrix metrics to a statsD backend
func WithStatsDCollector(addr, prefix string, sampleRate float32, flushBytes int) Option {
	return func(c *options) {
		c.statsD = &plugins.StatsdCollectorConfig{StatsdAddr: addr, Prefix: prefix, SampleRate: sampleRate, FlushBytes: flushBytes}
	}
}
