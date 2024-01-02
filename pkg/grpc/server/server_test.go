package server

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var fn = func(s *grpc.Server) {
	// pb.RegisterGreeterServer(s, &greeterServer{})
}

var unaryInterceptors = []grpc.UnaryServerInterceptor{
	func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		return nil, nil
	},
}

var streamInterceptors = []grpc.StreamServerInterceptor{
	func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return nil
	},
}

func TestRun(t *testing.T) {
	port := 50082
	Run(port, fn,
		WithSecure(insecure.NewCredentials()),
		WithUnaryInterceptor(unaryInterceptors...),
		WithStreamInterceptor(streamInterceptors...),
		WithServiceRegister(func() {}),
	)
	t.Log("grpc server started", port)
	time.Sleep(time.Second * 5)
}
