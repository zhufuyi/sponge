package interceptor

import (
	"google.golang.org/grpc"

	"github.com/zhufuyi/sponge/pkg/grpc/metrics"
)

// UnaryClientMetrics client-side metrics unary interceptor
func UnaryClientMetrics() grpc.UnaryClientInterceptor {
	return metrics.UnaryClientMetrics()
}

// StreamClientMetrics client-side metrics stream interceptor
func StreamClientMetrics() grpc.StreamClientInterceptor {
	return metrics.StreamClientMetrics()
}

// UnaryServerMetrics server-side metrics unary interceptor
func UnaryServerMetrics(opts ...metrics.Option) grpc.UnaryServerInterceptor {
	return metrics.UnaryServerMetrics(opts...)
}

// StreamServerMetrics server-side metrics stream interceptor
func StreamServerMetrics(opts ...metrics.Option) grpc.StreamServerInterceptor {
	return metrics.StreamServerMetrics(opts...)
}
