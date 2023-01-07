package grpccli

import (
	"context"
	"errors"

	"github.com/zhufuyi/sponge/pkg/grpc/interceptor"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/servicerd/discovery"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Dial secure connections
func Dial(ctx context.Context, endpoint string, opts ...Option) (*grpc.ClientConn, error) {
	return dial(ctx, endpoint, true, opts...)
}

// DialInsecure unsafe connections
func DialInsecure(ctx context.Context, endpoint string, opts ...Option) (*grpc.ClientConn, error) {
	return dial(ctx, endpoint, false, opts...)
}

func dial(ctx context.Context, endpoint string, isSecure bool, opts ...Option) (*grpc.ClientConn, error) {
	o := defaultOptions()
	o.apply(opts...)

	var unaryClientInterceptors []grpc.UnaryClientInterceptor
	var streamClientInterceptors []grpc.StreamClientInterceptor

	var clientOptions []grpc.DialOption

	// determining whether to use service discovery
	if o.discovery != nil {
		clientOptions = append(clientOptions, grpc.WithResolvers(
			discovery.NewBuilder(
				o.discovery,
				discovery.WithInsecure(!isSecure),
			)))
	}

	// secure connection or not
	if isSecure {
		if o.credentials == nil {
			return nil, errors.New("unset tls credentials")
		}
		clientOptions = append(clientOptions, grpc.WithTransportCredentials(o.credentials))
	} else {
		clientOptions = append(clientOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

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

	// load balance
	if o.enableLoadBalance {
		clientOptions = append(clientOptions, grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`))
	}

	// circuit breaker
	if o.enableCircuitBreaker {
		unaryClientInterceptors = append(unaryClientInterceptors, interceptor.UnaryClientCircuitBreaker())
	}

	// retry
	if o.enableRetry {
		unaryClientInterceptors = append(unaryClientInterceptors, interceptor.UnaryClientRetry())
	}

	if o.enableTrace {
		unaryClientInterceptors = append(unaryClientInterceptors, interceptor.UnaryClientTracing())
	}

	unaryClientInterceptors = append(unaryClientInterceptors, o.unaryInterceptors...)
	streamClientInterceptors = append(streamClientInterceptors, o.streamInterceptors...)

	o.dialOptions = append(o.dialOptions,
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(unaryClientInterceptors...),
		))
	o.dialOptions = append(o.dialOptions,
		grpc.WithStreamInterceptor(
			grpc_middleware.ChainStreamClient(streamClientInterceptors...),
		))

	clientOptions = append(clientOptions, o.dialOptions...)

	return grpc.DialContext(ctx, endpoint, clientOptions...)
}
