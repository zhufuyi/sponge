// Package server is generic grpc server-side.
package server

import (
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// RegisterFn register object
type RegisterFn func(srv *grpc.Server)

// ServiceRegisterFn service register
type ServiceRegisterFn func()

// Option set server option
type Option func(*options)

type options struct {
	credentials        credentials.TransportCredentials
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	serviceRegisterFn  ServiceRegisterFn
}

func defaultServerOptions() *options {
	return &options{}
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithSecure set secure
func WithSecure(credential credentials.TransportCredentials) Option {
	return func(o *options) {
		o.credentials = credential
	}
}

// WithUnaryInterceptor set unary interceptor
func WithUnaryInterceptor(interceptors ...grpc.UnaryServerInterceptor) Option {
	return func(o *options) {
		o.unaryInterceptors = interceptors
	}
}

// WithStreamInterceptor set stream interceptor
func WithStreamInterceptor(interceptors ...grpc.StreamServerInterceptor) Option {
	return func(o *options) {
		o.streamInterceptors = interceptors
	}
}

// WithServiceRegister set service register
func WithServiceRegister(fn ServiceRegisterFn) Option {
	return func(o *options) {
		o.serviceRegisterFn = fn
	}
}

func customInterceptorOptions(o *options) []grpc.ServerOption {
	var opts []grpc.ServerOption

	if o.credentials != nil {
		opts = append(opts, grpc.Creds(o.credentials))
	}

	if len(o.unaryInterceptors) > 0 {
		option := grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(o.unaryInterceptors...),
		)
		opts = append(opts, option)
	}
	if len(o.streamInterceptors) > 0 {
		option := grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(o.streamInterceptors...),
		)
		opts = append(opts, option)
	}

	return opts
}

// Run grpc server
func Run(port int, registerFns []RegisterFn, options ...Option) {
	o := defaultServerOptions()
	o.apply(options...)

	// listening on TCP port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// create a grpc server where interceptors can be injected
	srv := grpc.NewServer(customInterceptorOptions(o)...)

	// register object to the server
	for _, fn := range registerFns {
		fn(srv)
	}

	// register service to target
	if o.serviceRegisterFn != nil {
		o.serviceRegisterFn()
	}

	go func() {
		// run the server
		err = srv.Serve(listener)
		if err != nil {
			panic(err)
		}
	}()
}
