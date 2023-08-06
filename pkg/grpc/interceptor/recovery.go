package interceptor

import (
	"context"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ---------------------------------- client interceptor ----------------------------------

// UnaryClientRecovery client-side unary recovery
func UnaryClientRecovery() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = status.Errorf(codes.Internal, "triggered panic: %v", r)
			}
		}()

		err = invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}

// StreamClientRecovery client-side recovery stream interceptor
func StreamClientRecovery() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
		streamer grpc.Streamer, opts ...grpc.CallOption) (s grpc.ClientStream, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = status.Errorf(codes.Internal, "triggered panic: %v", r)
			}
		}()

		s, err = streamer(ctx, desc, cc, method, opts...)
		return s, err
	}
}

// ---------------------------------- server interceptor ----------------------------------

// UnaryServerRecovery recovery unary interceptor
func UnaryServerRecovery() grpc.UnaryServerInterceptor {
	// https://pkg.go.dev/github.com/grpc-ecosystem/go-grpc-middleware/recovery
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Internal, "triggered panic: %v", p)
	}
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}

	return grpc_recovery.UnaryServerInterceptor(opts...)
}

// StreamServerRecovery recovery stream interceptor
func StreamServerRecovery() grpc.StreamServerInterceptor {
	// https://pkg.go.dev/github.com/grpc-ecosystem/go-grpc-middleware/recovery
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Internal, "triggered panic: %v", p)
	}
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}

	return grpc_recovery.StreamServerInterceptor(opts...)
}
