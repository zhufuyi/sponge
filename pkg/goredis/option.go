package goredis

// Option set the redis options.
type Option func(*options)

type options struct {
	enableTrace bool
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// 默认设置
func defaultOptions() *options {
	return &options{
		enableTrace: false, // 是否开启链路跟踪，默认关闭
	}
}

// WithEnableTrace use trace
func WithEnableTrace() Option {
	return func(o *options) {
		o.enableTrace = true
	}
}
