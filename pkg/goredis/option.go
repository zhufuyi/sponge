package goredis

import (
	"crypto/tls"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Option set the redis options.
type Option func(*options)

type options struct {
	dialTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
	tlsConfig    *tls.Config

	// Note: this field is only used for Init and InitSingle, and the other parameters will be ignored.
	singleOptions *redis.Options

	// Note: this field is only used for InitSentinel, and the other parameters will be ignored.
	sentinelOptions *redis.FailoverOptions

	// Note: this field is only used for InitCluster, and the other parameters will be ignored.
	clusterOptions *redis.ClusterOptions

	// deprecated: use tp instead
	enableTrace    bool
	tracerProvider *trace.TracerProvider
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// default settings
func defaultOptions() *options {
	return &options{
		enableTrace: false, // whether to enable trace, default off
	}
}

// WithEnableTrace use trace, redis v8
// Deprecated: use WithEnableTracer instead
func WithEnableTrace() Option {
	return func(o *options) {
		o.enableTrace = true
	}
}

// WithTracing set redis tracer provider, redis v9
func WithTracing(tp *trace.TracerProvider) Option {
	return func(o *options) {
		o.tracerProvider = tp
	}
}

// WithDialTimeout set dail timeout
func WithDialTimeout(t time.Duration) Option {
	return func(o *options) {
		o.dialTimeout = t
	}
}

// WithReadTimeout set read timeout
func WithReadTimeout(t time.Duration) Option {
	return func(o *options) {
		o.readTimeout = t
	}
}

// WithWriteTimeout set write timeout
func WithWriteTimeout(t time.Duration) Option {
	return func(o *options) {
		o.writeTimeout = t
	}
}

// WithTLSConfig set TLS config
func WithTLSConfig(c *tls.Config) Option {
	return func(o *options) {
		o.tlsConfig = c
	}
}

// WithSingleOptions set single redis options
func WithSingleOptions(opt *redis.Options) Option {
	return func(o *options) {
		o.singleOptions = opt
	}
}

// WithSentinelOptions set redis sentinel options
func WithSentinelOptions(opt *redis.FailoverOptions) Option {
	return func(o *options) {
		o.sentinelOptions = opt
	}
}

// WithClusterOptions set redis cluster options
func WithClusterOptions(opt *redis.ClusterOptions) Option {
	return func(o *options) {
		o.clusterOptions = opt
	}
}
