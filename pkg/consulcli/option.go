package consulcli

import (
	"time"
)

// Option set the consul client options.
type Option func(*options)

type options struct {
	scheme     string
	waitTime   time.Duration
	datacenter string
}

func defaultOptions() *options {
	return &options{
		scheme:   "http",
		waitTime: time.Second * 5,
	}
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithWaitTime set wait time
func WithWaitTime(waitTime time.Duration) Option {
	return func(o *options) {
		o.waitTime = waitTime
	}
}

// WithScheme set scheme
func WithScheme(scheme string) Option {
	return func(o *options) {
		o.scheme = scheme
	}
}

// WithDatacenter set datacenter
func WithDatacenter(datacenter string) Option {
	return func(o *options) {
		o.datacenter = datacenter
	}
}
