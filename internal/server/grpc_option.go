package server

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
)

// GrpcOption grpc settings
type GrpcOption func(*grpcOptions)

type grpcOptions struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
	instance     *registry.ServiceInstance
	iRegistry    registry.Registry
}

func defaultGrpcOptions() *grpcOptions {
	return &grpcOptions{
		readTimeout:  time.Second * 3,
		writeTimeout: time.Second * 3,
		instance:     nil,
		iRegistry:    nil,
	}
}

func (o *grpcOptions) apply(opts ...GrpcOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithGrpcReadTimeout setting up read timeout
func WithGrpcReadTimeout(timeout time.Duration) GrpcOption {
	return func(o *grpcOptions) {
		o.readTimeout = timeout
	}
}

// WithGrpcWriteTimeout setting up writer timeout
func WithGrpcWriteTimeout(timeout time.Duration) GrpcOption {
	return func(o *grpcOptions) {
		o.writeTimeout = timeout
	}
}

// WithGrpcRegistry registration services
func WithGrpcRegistry(iRegistry registry.Registry, instance *registry.ServiceInstance) GrpcOption {
	return func(o *grpcOptions) {
		o.iRegistry = iRegistry
		o.instance = instance
	}
}
