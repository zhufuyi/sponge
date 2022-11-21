package grpccli

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/servicerd/registry"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Option grpc dial options
type Option func(*options)

// options grpc dial options
type options struct {
	timeout time.Duration

	credentials        credentials.TransportCredentials // secure connections
	dialOptions        []grpc.DialOption                // custom options
	unaryInterceptors  []grpc.UnaryClientInterceptor    // custom unary interceptor
	streamInterceptors []grpc.StreamClientInterceptor   // custom stream interceptor

	enableLog bool // whether to turn on the log
	log       *zap.Logger

	enableTrace          bool // whether to turn on tracing
	enableMetrics        bool // whether to turn on metrics
	enableRetry          bool // whether to turn on retry
	enableLoadBalance    bool // whether to turn on load balance
	enableCircuitBreaker bool // whether to turn on circuit breaker

	discovery registry.Discovery // if not nil means use service discovery
}

func defaultOptions() *options {
	return &options{
		enableLog: false,

		timeout:            time.Second * 5,
		credentials:        nil,
		dialOptions:        nil,
		unaryInterceptors:  nil,
		streamInterceptors: nil,
		discovery:          nil,
	}
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithTimeout set dial timeout
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

// WithEnableLog enable log
func WithEnableLog(log *zap.Logger) Option {
	return func(o *options) {
		o.enableLog = true
		o.log = log
	}
}

// WithEnableTrace enable trace
func WithEnableTrace() Option {
	return func(o *options) {
		o.enableTrace = true
	}
}

// WithEnableMetrics enable metrics
func WithEnableMetrics() Option {
	return func(o *options) {
		o.enableMetrics = true
	}
}

// WithEnableLoadBalance enable load balance
func WithEnableLoadBalance() Option {
	return func(o *options) {
		o.enableLoadBalance = true
	}
}

// WithEnableRetry enable registry
func WithEnableRetry() Option {
	return func(o *options) {
		o.enableRetry = true
	}
}

// WithEnableCircuitBreaker enable circuit breaker
func WithEnableCircuitBreaker() Option {
	return func(o *options) {
		o.enableCircuitBreaker = true
	}
}

// WithCredentials set dial credentials
func WithCredentials(credentials credentials.TransportCredentials) Option {
	return func(o *options) {
		o.credentials = credentials
	}
}

// WithDialOptions set dial options
func WithDialOptions(dialOptions ...grpc.DialOption) Option {
	return func(o *options) {
		o.dialOptions = append(o.dialOptions, dialOptions...)
	}
}

// WithUnaryInterceptors set dial unaryInterceptors
func WithUnaryInterceptors(unaryInterceptors ...grpc.UnaryClientInterceptor) Option {
	return func(o *options) {
		o.unaryInterceptors = append(o.unaryInterceptors, unaryInterceptors...)
	}
}

// WithStreamInterceptors set dial streamInterceptors
func WithStreamInterceptors(streamInterceptors ...grpc.StreamClientInterceptor) Option {
	return func(o *options) {
		o.streamInterceptors = append(o.streamInterceptors, streamInterceptors...)
	}
}

// WithDiscovery set dial discovery
func WithDiscovery(discovery registry.Discovery) Option {
	return func(o *options) {
		o.discovery = discovery
	}
}
