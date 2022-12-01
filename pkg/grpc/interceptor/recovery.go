package interceptor

import (
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
