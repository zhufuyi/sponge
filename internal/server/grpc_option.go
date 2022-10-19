package server

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
)

// GRPCOption grpc settings
type GRPCOption func(*grpcOptions)

type grpcOptions struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
	instance     *registry.ServiceInstance
	iRegistry    registry.Registry
}

func defaultGRPCOptions() *grpcOptions {
	return &grpcOptions{
		readTimeout:  time.Second * 3,
		writeTimeout: time.Second * 3,
		instance:     nil,
		iRegistry:    nil,
	}
}

func (o *grpcOptions) apply(opts ...GRPCOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithGRPCReadTimeout setting up read timeout
func WithGRPCReadTimeout(timeout time.Duration) GRPCOption {
	return func(o *grpcOptions) {
		o.readTimeout = timeout
	}
}

// WithGRPCWriteTimeout setting up writer timeout
func WithGRPCWriteTimeout(timeout time.Duration) GRPCOption {
	return func(o *grpcOptions) {
		o.writeTimeout = timeout
	}
}

// WithRegistry setting up registry
func WithRegistry(iRegistry registry.Registry, instance *registry.ServiceInstance) GRPCOption {
	return func(o *grpcOptions) {
		o.iRegistry = iRegistry
		o.instance = instance
	}
}
