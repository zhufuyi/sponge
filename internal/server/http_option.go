package server

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
)

// HTTPOption setting up http
type HTTPOption func(*httpOptions)

type httpOptions struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
	isProd       bool

	instance  *registry.ServiceInstance
	iRegistry registry.Registry
}

func defaultHTTPOptions() *httpOptions {
	return &httpOptions{
		readTimeout:  time.Second * 60,
		writeTimeout: time.Second * 60,
		isProd:       false,
		instance:     nil,
		iRegistry:    nil,
	}
}

func (o *httpOptions) apply(opts ...HTTPOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithHTTPReadTimeout setting up read timeout
func WithHTTPReadTimeout(timeout time.Duration) HTTPOption {
	return func(o *httpOptions) {
		o.readTimeout = timeout
	}
}

// WithHTTPWriteTimeout setting up writer timeout
func WithHTTPWriteTimeout(timeout time.Duration) HTTPOption {
	return func(o *httpOptions) {
		o.writeTimeout = timeout
	}
}

// WithHTTPIsProd setting up production environment markers
func WithHTTPIsProd(isProd bool) HTTPOption {
	return func(o *httpOptions) {
		o.isProd = isProd
	}
}

// WithHTTPRegistry registration services
func WithHTTPRegistry(iRegistry registry.Registry, instance *registry.ServiceInstance) HTTPOption {
	return func(o *httpOptions) {
		o.iRegistry = iRegistry
		o.instance = instance
	}
}
