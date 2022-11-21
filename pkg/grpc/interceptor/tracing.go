package interceptor

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// UnaryClientTracing client-side tracing unary interceptor
func UnaryClientTracing() grpc.UnaryClientInterceptor {
	return otelgrpc.UnaryClientInterceptor()
}

// StreamClientTracing client-side tracing stream interceptor
func StreamClientTracing() grpc.StreamClientInterceptor {
	return otelgrpc.StreamClientInterceptor()
}

// UnaryServerTracing server-side tracing unary interceptor
func UnaryServerTracing() grpc.UnaryServerInterceptor {
	return otelgrpc.UnaryServerInterceptor()
}

// StreamServerTracing server-side tracing stream interceptor
func StreamServerTracing() grpc.StreamServerInterceptor {
	return otelgrpc.StreamServerInterceptor()
}
