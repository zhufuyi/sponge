package ratelimiter

import (
	"golang.org/x/time/rate"
)

var (
	// default qps value
	defaultQPS rate.Limit = 500

	// default the maximum instantaneous request spike allowed, burst >= qps
	defaultBurst = 1000

	// default is path limit, fault:path limit, true:ip limit
	defaultIsIP = false //nolint
)

// Option set the rate limits options.
type Option func(*options)

func defaultOptions() *options {
	return &options{
		qps:   defaultQPS,
		burst: defaultBurst,
		isIP:  false,
	}
}

type options struct {
	qps   rate.Limit
	burst int
	isIP  bool // false: path limit, true: IP limit
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithQPS set the qps value
func WithQPS(qps int) Option {
	return func(o *options) {
		o.qps = rate.Limit(qps)
	}
}

// WithBurst set the burst value, burst >= qps
func WithBurst(burst int) Option {
	return func(o *options) {
		o.burst = burst
	}
}

// WithPath set the path limit mode
func WithPath() Option {
	return func(o *options) {
		o.isIP = false
	}
}

// WithIP set the path limit mode
func WithIP() Option {
	return func(o *options) {
		o.isIP = true
	}
}
