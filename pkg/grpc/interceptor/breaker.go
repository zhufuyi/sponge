package interceptor

import (
	"context"
	"github.com/zhufuyi/sponge/pkg/container/group"
	"github.com/zhufuyi/sponge/pkg/errcode"

	"github.com/go-kratos/aegis/circuitbreaker"
	"github.com/go-kratos/aegis/circuitbreaker/sre"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrNotAllowed error not allowed.
var ErrNotAllowed = circuitbreaker.ErrNotAllowed

// CircuitBreakerOption set the circuit breaker circuitBreakerOptions.
type CircuitBreakerOption func(*circuitBreakerOptions)

type circuitBreakerOptions struct {
	group *group.Group
}

func defaultCircuitBreakerOptions() *circuitBreakerOptions {
	return &circuitBreakerOptions{
		group: group.NewGroup(func() interface{} {
			return sre.NewBreaker()
		}),
	}
}

func (o *circuitBreakerOptions) apply(opts ...CircuitBreakerOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithGroup with circuit breaker group.
// NOTE: implements generics circuitbreaker.CircuitBreaker
func WithGroup(g *group.Group) CircuitBreakerOption {
	return func(o *circuitBreakerOptions) {
		o.group = g
	}
}

// UnaryClientCircuitBreaker client-side unary circuit breaker interceptor
func UnaryClientCircuitBreaker(opts ...CircuitBreakerOption) grpc.UnaryClientInterceptor {
	o := defaultCircuitBreakerOptions()
	o.apply(opts...)

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		breaker := o.group.Get(method).(circuitbreaker.CircuitBreaker)
		if err := breaker.Allow(); err != nil {
			// NOTE: when client reject request locally,
			// continue add counter let the drop ratio higher.
			breaker.MarkFailed()
			return errcode.StatusServiceUnavailable.ToRPCErr(err.Error())
		}

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			// NOTE: need to check internal and service unavailable error
			s, ok := status.FromError(err)
			if ok && (s.Code() == codes.Internal || s.Code() == codes.Unavailable) {
				breaker.MarkFailed()
			} else {
				breaker.MarkSuccess()
			}
		}

		return err
	}
}

// SteamClientCircuitBreaker client-side stream circuit breaker interceptor
func SteamClientCircuitBreaker(opts ...CircuitBreakerOption) grpc.StreamClientInterceptor {
	o := defaultCircuitBreakerOptions()
	o.apply(opts...)

	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		breaker := o.group.Get(method).(circuitbreaker.CircuitBreaker)
		if err := breaker.Allow(); err != nil {
			// NOTE: when client reject request locally,
			// continue add counter let the drop ratio higher.
			breaker.MarkFailed()
			return nil, errcode.StatusServiceUnavailable.ToRPCErr(err.Error())
		}

		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			// NOTE: need to check internal and service unavailable error
			s, ok := status.FromError(err)
			if ok && (s.Code() == codes.Internal || s.Code() == codes.Unavailable) {
				breaker.MarkFailed()
			} else {
				breaker.MarkSuccess()
			}
		}

		return clientStream, err
	}
}

// UnaryServerCircuitBreaker server-side unary circuit breaker interceptor
func UnaryServerCircuitBreaker(opts ...CircuitBreakerOption) grpc.UnaryServerInterceptor {
	o := defaultCircuitBreakerOptions()
	o.apply(opts...)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		breaker := o.group.Get(info.FullMethod).(circuitbreaker.CircuitBreaker)
		if err := breaker.Allow(); err != nil {
			// NOTE: when client reject request locally,
			// continue add counter let the drop ratio higher.
			breaker.MarkFailed()
			return nil, errcode.StatusServiceUnavailable.ToRPCErr(err.Error())
		}

		reply, err := handler(ctx, req)
		if err != nil {
			// NOTE: need to check internal and service unavailable error
			s, ok := status.FromError(err)
			if ok && (s.Code() == codes.Internal || s.Code() == codes.Unavailable) {
				breaker.MarkFailed()
			} else {
				breaker.MarkSuccess()
			}
		}

		return reply, err
	}
}

// SteamServerCircuitBreaker server-side stream circuit breaker interceptor
func SteamServerCircuitBreaker(opts ...CircuitBreakerOption) grpc.StreamServerInterceptor {
	o := defaultCircuitBreakerOptions()
	o.apply(opts...)

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		breaker := o.group.Get(info.FullMethod).(circuitbreaker.CircuitBreaker)
		if err := breaker.Allow(); err != nil {
			// NOTE: when client reject request locally,
			// continue add counter let the drop ratio higher.
			breaker.MarkFailed()
			return errcode.StatusServiceUnavailable.ToRPCErr(err.Error())
		}

		err := handler(srv, ss)
		if err != nil {
			// NOTE: need to check internal and service unavailable error
			s, ok := status.FromError(err)
			if ok && (s.Code() == codes.Internal || s.Code() == codes.Unavailable) {
				breaker.MarkFailed()
			} else {
				breaker.MarkSuccess()
			}
		}

		return err
	}
}
