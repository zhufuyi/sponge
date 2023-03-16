package grpccli

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/servicerd/registry"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	secureOneWay = "one-way"
	secureTwoWay = "two-way"
)

// Option grpc dial options
type Option func(*options)

// options grpc dial options
type options struct {
	timeout time.Duration

	// secure setting
	secureType string // secure type "","one-way","two-way"
	serverName string // server name
	caFile     string // ca file
	certFile   string // cert file
	keyFile    string // key file

	// token setting
	enableToken bool // whether to turn on token
	appID       string
	appKey      string

	// interceptor setting
	enableLog            bool // whether to turn on the log
	log                  *zap.Logger
	enableRequestID      bool               // whether to turn on the request id
	enableTrace          bool               // whether to turn on tracing
	enableMetrics        bool               // whether to turn on metrics
	enableRetry          bool               // whether to turn on retry
	enableLoadBalance    bool               // whether to turn on load balance
	enableCircuitBreaker bool               // whether to turn on circuit breaker
	discovery            registry.Discovery // if not nil means use service discovery

	// custom setting
	dialOptions        []grpc.DialOption              // custom options
	unaryInterceptors  []grpc.UnaryClientInterceptor  // custom unary interceptor
	streamInterceptors []grpc.StreamClientInterceptor // custom stream interceptor
}

func defaultOptions() *options {
	return &options{
		secureType: "",
		serverName: "localhost",
		certFile:   "",
		keyFile:    "",
		caFile:     "",

		enableLog: false,

		timeout:            time.Second * 5,
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

// WithEnableRequestID enable request id
func WithEnableRequestID() Option {
	return func(o *options) {
		o.enableRequestID = true
	}
}

// WithEnableLog enable log
func WithEnableLog(log *zap.Logger) Option {
	return func(o *options) {
		o.enableLog = true
		if log != nil {
			o.log = log
		} else {
			o.log = zap.NewNop()
		}
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

func (o *options) isSecure() bool {
	if o.secureType == secureOneWay || o.secureType == secureTwoWay {
		return true
	}
	return false
}

// WithSecure support setting one-way or two-way secure
func WithSecure(t string, serverName string, caFile string, certFile string, keyFile string) Option {
	switch t {
	case secureOneWay:
		return WithOneWaySecure(serverName, certFile)
	case secureTwoWay:
		return WithTwoWaySecure(serverName, caFile, certFile, keyFile)
	}

	return func(o *options) {
		o.secureType = t
	}
}

// WithOneWaySecure set one-way secure
func WithOneWaySecure(serverName string, certFile string) Option {
	return func(o *options) {
		if serverName == "" {
			serverName = "localhost"
		}
		o.secureType = secureOneWay
		o.serverName = serverName
		o.certFile = certFile
	}
}

// WithTwoWaySecure set two-way secure
func WithTwoWaySecure(serverName string, caFile string, certFile string, keyFile string) Option {
	return func(o *options) {
		if serverName == "" {
			serverName = "localhost"
		}
		o.secureType = secureTwoWay
		o.serverName = serverName
		o.caFile = caFile
		o.certFile = certFile
		o.keyFile = keyFile
	}
}

// WithToken set token
func WithToken(enable bool, appID string, appKey string) Option {
	return func(o *options) {
		o.enableToken = enable
		o.appID = appID
		o.appKey = appKey
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
