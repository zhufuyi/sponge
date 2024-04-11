package server

import (
	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
)

// HTTPOption setting up http
type HTTPOption func(*httpOptions)

type httpOptions struct {
	isProd    bool
	instance  *registry.ServiceInstance
	iRegistry registry.Registry
}

func defaultHTTPOptions() *httpOptions {
	return &httpOptions{
		isProd:    false,
		instance:  nil,
		iRegistry: nil,
	}
}

func (o *httpOptions) apply(opts ...HTTPOption) {
	for _, opt := range opts {
		opt(o)
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
