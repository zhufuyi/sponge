package goredis

import (
	"crypto/tls"
	"time"
)

// Option set the redis options.
type Option func(*options)

type options struct {
	enableTrace  bool
	dialTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration

	tlsConfig *tls.Config
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

// WithTLSConfig set TLS config
func WithTLSConfig(c *tls.Config) Option {
	return func(o *options) {
		o.tlsConfig = c
	}
}
