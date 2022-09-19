package interceptor

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// UnaryClientTracing 客户端链路跟踪unary拦截器
func UnaryClientTracing() grpc.UnaryClientInterceptor {
	return otelgrpc.UnaryClientInterceptor()
}

// StreamClientTracing 客户端链路跟踪stream拦截器
func StreamClientTracing() grpc.StreamClientInterceptor {
	return otelgrpc.StreamClientInterceptor()
}

// UnaryServerTracing 服务端链路跟踪unary拦截器
func UnaryServerTracing() grpc.UnaryServerInterceptor {
	return otelgrpc.UnaryServerInterceptor()
}

// StreamServerTracing 服务端链路跟踪stream拦截器
func StreamServerTracing() grpc.StreamServerInterceptor {
	return otelgrpc.StreamServerInterceptor()
}
