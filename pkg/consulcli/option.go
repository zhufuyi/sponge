package consulcli

import (
	"time"

	"github.com/hashicorp/consul/api"
)

// Option set the consul client options.
type Option func(*options)

type options struct {
	scheme     string
	waitTime   time.Duration
	datacenter string

	// if you set this parameter, all fields above are invalid
	config *api.Config
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

// WithConfig set consul config
func WithConfig(c *api.Config) Option {
	return func(o *options) {
		o.config = c
	}
}
