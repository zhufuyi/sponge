package interceptor

import (
	"github.com/zhufuyi/sponge/pkg/grpc/metrics"

	"google.golang.org/grpc"
)

// UnaryClientMetrics 客户端指标unary拦截器
func UnaryClientMetrics() grpc.UnaryClientInterceptor {
	return metrics.UnaryClientMetrics()
}

// StreamClientMetrics 客户端指标stream拦截器
func StreamClientMetrics() grpc.StreamClientInterceptor {
	return metrics.StreamClientMetrics()
}

// UnaryServerMetrics 服务端指标unary拦截器
func UnaryServerMetrics(opts ...metrics.Option) grpc.UnaryServerInterceptor {
	return metrics.UnaryServerMetrics(opts...)
}

// StreamServerMetrics 服务端指标stream拦截器
func StreamServerMetrics(opts ...metrics.Option) grpc.StreamServerInterceptor {
	return metrics.StreamServerMetrics(opts...)
}
