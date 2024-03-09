package interceptor

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// UnaryClientTracing client-side tracing unary interceptor
func UnaryClientTracing() grpc.UnaryClientInterceptor {
	return otelgrpc.UnaryClientInterceptor() //nolint
}

// StreamClientTracing client-side tracing stream interceptor
func StreamClientTracing() grpc.StreamClientInterceptor {
	return otelgrpc.StreamClientInterceptor() //nolint
}

// UnaryServerTracing server-side tracing unary interceptor
func UnaryServerTracing() grpc.UnaryServerInterceptor {
	return otelgrpc.UnaryServerInterceptor() //nolint
}

// StreamServerTracing server-side tracing stream interceptor
func StreamServerTracing() grpc.StreamServerInterceptor {
	return otelgrpc.StreamServerInterceptor() //nolint
}
