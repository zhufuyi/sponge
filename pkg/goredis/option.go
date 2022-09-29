package goredis

import "time"

// Option set the redis options.
type Option func(*options)

type options struct {
	enableTrace  bool
	dialTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
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
