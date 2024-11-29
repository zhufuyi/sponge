package client

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

type builder struct{}

func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	return nil, nil
}

func (b *builder) Scheme() string {
	return ""
}

var unaryInterceptors = []grpc.UnaryClientInterceptor{
	func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return nil
	},
}

var streamInterceptors = []grpc.StreamClientInterceptor{
	func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return nil, nil
	},
}

func TestNewClient(t *testing.T) {
	conn, err := NewClient("127.0.0.1:50082",
		WithServiceDiscover(new(builder)),
		WithLoadBalance(),
		WithSecure(insecure.NewCredentials()),
		WithUnaryInterceptor(unaryInterceptors...),
		WithStreamInterceptor(streamInterceptors...),
		WithDialOption(grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`)),
	)
	defer conn.Close()
	t.Log(conn, err)
	time.Sleep(time.Second)
}
