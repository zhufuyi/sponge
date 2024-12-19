package server

import (
	"context"
	"fmt"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-dev-frame/sponge/pkg/grpc/metrics"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/utils"
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
	port, _ := utils.GetAvailablePort()
	Run(port, fn,
		WithSecure(insecure.NewCredentials()),
		WithUnaryInterceptor(unaryInterceptors...),
		WithStreamInterceptor(streamInterceptors...),
		WithServiceRegister(func() {}),
		WithStatConnections(metrics.WithConnectionsLogger(logger.Get()), metrics.WithConnectionsGauge()),
	)
	t.Log("grpc server started", port)
	time.Sleep(time.Second * 2)

	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(time.Second * 2)
	_ = conn.Close()
	time.Sleep(time.Second * 1)
}
