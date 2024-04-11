package server

import (
	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
)

// GrpcOption grpc settings
type GrpcOption func(*grpcOptions)

type grpcOptions struct {
	instance  *registry.ServiceInstance
	iRegistry registry.Registry
}

func defaultGrpcOptions() *grpcOptions {
	return &grpcOptions{
		instance:  nil,
		iRegistry: nil,
	}
}

func (o *grpcOptions) apply(opts ...GrpcOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithGrpcRegistry registration services
func WithGrpcRegistry(iRegistry registry.Registry, instance *registry.ServiceInstance) GrpcOption {
	return func(o *grpcOptions) {
		o.iRegistry = iRegistry
		o.instance = instance
	}
}
