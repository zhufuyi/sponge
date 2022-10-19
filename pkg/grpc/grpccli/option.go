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

	credentials        credentials.TransportCredentials // 安全连接credentials
	dialOptions        []grpc.DialOption                // 自定义options
	unaryInterceptors  []grpc.UnaryClientInterceptor    // 自定义unary拦截器
	streamInterceptors []grpc.StreamClientInterceptor   // 自定义stream拦截器

	enableLog bool // 是否开启日志
	log       *zap.Logger

	enableTrace          bool // 是否开启链路跟踪
	enableMetrics        bool // 是否开启指标
	enableRetry          bool // 是否开启重试
	enableLoadBalance    bool // 是否开启负载均衡器
	enableCircuitBreaker bool // 是否开启熔断器

	discovery registry.Discovery // 服务发现接口
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
