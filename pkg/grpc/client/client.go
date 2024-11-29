// Package client is generic grpc client-side.
package client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

// Option client option func
type Option func(*options)

type options struct {
	builders           []resolver.Builder
	isLoadBalance      bool
	credentials        credentials.TransportCredentials
	unaryInterceptors  []grpc.UnaryClientInterceptor
	streamInterceptors []grpc.StreamClientInterceptor
	dialOptions        []grpc.DialOption
}

func defaultOptions() *options {
	return &options{}
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithServiceDiscover set service discover
func WithServiceDiscover(builders ...resolver.Builder) Option {
	return func(o *options) {
		o.builders = builders
	}
}

// WithLoadBalance set load balance
func WithLoadBalance() Option {
	return func(o *options) {
		o.isLoadBalance = true
	}
}

// WithSecure set secure
func WithSecure(credential credentials.TransportCredentials) Option {
	return func(o *options) {
		o.credentials = credential
	}
}

// WithUnaryInterceptor set unary interceptor
func WithUnaryInterceptor(interceptors ...grpc.UnaryClientInterceptor) Option {
	return func(o *options) {
		o.unaryInterceptors = interceptors
	}
}

// WithStreamInterceptor set stream interceptor
func WithStreamInterceptor(interceptors ...grpc.StreamClientInterceptor) Option {
	return func(o *options) {
		o.streamInterceptors = interceptors
	}
}

// WithDialOption set DialOption
func WithDialOption(dialOptions ...grpc.DialOption) Option {
	return func(o *options) {
		o.dialOptions = dialOptions
	}
}

// NewClient create a new grpc client
func NewClient(endpoint string, opts ...Option) (*grpc.ClientConn, error) {
	o := defaultOptions()
	o.apply(opts...)

	var dialOptions []grpc.DialOption

	// service discovery
	if len(o.builders) > 0 {
		dialOptions = append(dialOptions, grpc.WithResolvers(o.builders...))
	}

	// load balance option
	if o.isLoadBalance {
		dialOptions = append(dialOptions, grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`))
	}

	// secure option
	if o.credentials == nil {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(o.credentials))
	}

	// custom dial option
	if len(o.dialOptions) > 0 {
		dialOptions = append(dialOptions, o.dialOptions...)
	}

	// custom unary interceptor option
	if len(o.unaryInterceptors) > 0 {
		dialOptions = append(dialOptions, grpc.WithChainUnaryInterceptor(o.unaryInterceptors...))
	}

	// custom stream interceptor option
	if len(o.streamInterceptors) > 0 {
		dialOptions = append(dialOptions, grpc.WithChainStreamInterceptor(o.streamInterceptors...))
	}

	return grpc.NewClient(endpoint, dialOptions...)
}

// Dial to grpc server
func Dial(_ context.Context, endpoint string, opts ...Option) (*grpc.ClientConn, error) {
	return NewClient(endpoint, opts...)
}
