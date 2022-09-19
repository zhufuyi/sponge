package interceptor

import (
	"github.com/zhufuyi/sponge/pkg/grpc/hystrix"

	"google.golang.org/grpc"
)

// UnaryClientHystrix 客户端熔断器unary拦截器
func UnaryClientHystrix(commandName string, opts ...hystrix.Option) grpc.UnaryClientInterceptor {
	return hystrix.UnaryClientInterceptor(commandName, opts...)
}

// SteamClientHystrix 客户端熔断器stream拦截器
func SteamClientHystrix(commandName string, opts ...hystrix.Option) grpc.StreamClientInterceptor {
	return hystrix.StreamClientInterceptor(commandName, opts...)
}
