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
func WithSecure(credentials credentials.TransportCredentials) Option {
	return func(o *options) {
		o.credentials = credentials
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

// Dial to grpc server
func Dial(ctx context.Context, endpoint string, opts ...Option) (*grpc.ClientConn, error) {
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

	// custom unary interceptor option
	if len(o.unaryInterceptors) > 0 {
		dialOptions = append(dialOptions, grpc.WithChainUnaryInterceptor(o.unaryInterceptors...))
	}

	// custom stream interceptor option
	if len(o.streamInterceptors) > 0 {
		dialOptions = append(dialOptions, grpc.WithChainStreamInterceptor(o.streamInterceptors...))
	}

	return grpc.DialContext(ctx, endpoint, dialOptions...)
}
