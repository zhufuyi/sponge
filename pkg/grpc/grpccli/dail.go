// Package grpccli is grpc client with support for service discovery, logging, load balancing, trace, metrics, retries, circuit breaker.
package grpccli

import (
	"context"
	"errors"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zhufuyi/sponge/pkg/grpc/gtls"
	"github.com/zhufuyi/sponge/pkg/grpc/interceptor"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/servicerd/discovery"
)

// NewClient creates a new grpc client
func NewClient(endpoint string, opts ...Option) (*grpc.ClientConn, error) {
	o := defaultOptions()
	o.apply(opts...)

	var clientOptions []grpc.DialOption

	// service discovery
	if o.discovery != nil {
		clientOptions = append(clientOptions, grpc.WithResolvers(
			discovery.NewBuilder(
				o.discovery,
				discovery.WithInsecure(o.discoveryInsecure),
			)))
	}

	// load balance option
	if o.enableLoadBalance {
		clientOptions = append(clientOptions, grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`))
	}

	// secure option
	so, err := secureOption(o)
	if err != nil {
		return nil, err
	}
	clientOptions = append(clientOptions, so)

	// token option
	if o.enableToken {
		clientOptions = append(clientOptions, interceptor.ClientTokenOption(
			o.appID,
			o.appKey,
			o.isSecure(),
		))
	}

	// unary options
	clientOptions = append(clientOptions, unaryClientOptions(o))
	// stream options
	clientOptions = append(clientOptions, streamClientOptions(o))
	// custom options
	clientOptions = append(clientOptions, o.dialOptions...)

	return grpc.NewClient(endpoint, clientOptions...)
}

// Dial to grpc server
// Deprecated: use NewClient instead
func Dial(_ context.Context, endpoint string, opts ...Option) (*grpc.ClientConn, error) {
	return NewClient(endpoint, opts...)
}

func secureOption(o *options) (grpc.DialOption, error) {
	switch o.secureType {
	case secureOneWay: // server side certification
		if o.certFile == "" {
			return nil, errors.New("cert file is empty")
		}
		credentials, err := gtls.GetClientTLSCredentials(o.serverName, o.certFile)
		if err != nil {
			return nil, err
		}
		return grpc.WithTransportCredentials(credentials), nil

	case secureTwoWay: // both client and server side certification
		if o.caFile == "" {
			return nil, errors.New("ca file is empty")
		}
		if o.certFile == "" {
			return nil, errors.New("cert file is empty")
		}
		if o.keyFile == "" {
			return nil, errors.New("key file is empty")
		}
		credentials, err := gtls.GetClientTLSCredentialsByCA(
			o.serverName,
			o.caFile,
			o.certFile,
			o.keyFile,
		)
		if err != nil {
			return nil, err
		}
		return grpc.WithTransportCredentials(credentials), nil

	default:
		return grpc.WithTransportCredentials(insecure.NewCredentials()), nil
	}
}

func unaryClientOptions(o *options) grpc.DialOption {
	var unaryClientInterceptors []grpc.UnaryClientInterceptor

	unaryClientInterceptors = append(unaryClientInterceptors, interceptor.UnaryClientRecovery())

	if o.requestTimeout > 0 {
		unaryClientInterceptors = append(unaryClientInterceptors, interceptor.UnaryClientTimeout(o.requestTimeout))
	}

	// request id
	if o.enableRequestID {
		unaryClientInterceptors = append(unaryClientInterceptors, interceptor.UnaryClientRequestID())
	}

	// logging
	if o.enableLog {
		unaryClientInterceptors = append(unaryClientInterceptors, interceptor.UnaryClientLog(logger.Get()))
	}

	// metrics
	if o.enableMetrics {
		unaryClientInterceptors = append(unaryClientInterceptors, interceptor.UnaryClientMetrics())
	}

	// circuit breaker
	//if o.enableCircuitBreaker {
	//	unaryClientInterceptors = append(unaryClientInterceptors, interceptor.UnaryClientCircuitBreaker(
	//	// set rpc code for circuit breaker, default already includes codes.Internal and codes.Unavailable
	//	//interceptor.WithValidCode(codes.PermissionDenied),
	//	))
	//}

	// retry
	if o.enableRetry {
		unaryClientInterceptors = append(unaryClientInterceptors, interceptor.UnaryClientRetry())
	}

	// trace
	if o.enableTrace {
		unaryClientInterceptors = append(unaryClientInterceptors, interceptor.UnaryClientTracing())
	}

	// custom unary interceptors
	unaryClientInterceptors = append(unaryClientInterceptors, o.unaryInterceptors...)

	return grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(unaryClientInterceptors...),
	)
}

func streamClientOptions(o *options) grpc.DialOption {
	var streamClientInterceptors []grpc.StreamClientInterceptor

	streamClientInterceptors = append(streamClientInterceptors, interceptor.StreamClientRecovery())

	// request id
	if o.enableRequestID {
		streamClientInterceptors = append(streamClientInterceptors, interceptor.StreamClientRequestID())
	}

	// logging
	if o.enableLog {
		streamClientInterceptors = append(streamClientInterceptors, interceptor.StreamClientLog(logger.Get()))
	}

	// metrics
	if o.enableMetrics {
		streamClientInterceptors = append(streamClientInterceptors, interceptor.StreamClientMetrics())
	}

	// circuit breaker
	//if o.enableCircuitBreaker {
	//	streamClientInterceptors = append(streamClientInterceptors, interceptor.StreamClientCircuitBreaker(
	//	// set rpc code for circuit breaker, default already includes codes.Internal and codes.Unavailable
	//	//interceptor.WithValidCode(codes.PermissionDenied),
	//	))
	//}

	// retry
	if o.enableRetry {
		streamClientInterceptors = append(streamClientInterceptors, interceptor.StreamClientRetry())
	}

	// trace
	if o.enableTrace {
		streamClientInterceptors = append(streamClientInterceptors, interceptor.StreamClientTracing())
	}

	// custom stream interceptors
	streamClientInterceptors = append(streamClientInterceptors, o.streamInterceptors...)

	return grpc.WithStreamInterceptor(
		grpc_middleware.ChainStreamClient(streamClientInterceptors...),
	)
}
