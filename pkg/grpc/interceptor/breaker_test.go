package interceptor

import (
	"context"
	"testing"

	"github.com/zhufuyi/sponge/pkg/container/group"
	"github.com/zhufuyi/sponge/pkg/errcode"

	"github.com/go-kratos/aegis/circuitbreaker/sre"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestUnaryClientCircuitBreaker(t *testing.T) {
	interceptor := UnaryClientCircuitBreaker(WithGroup(
		group.NewGroup(func() interface{} {
			return sre.NewBreaker()
		}),
	))
	assert.NotNil(t, interceptor)

	ivoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return errcode.StatusInternalServerError.ToRPCErr()
	}
	for i := 0; i < 110; i++ {
		err := interceptor(context.Background(), "/test", nil, nil, nil, ivoker)
		assert.Error(t, err)
	}

	ivoker = func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return errcode.StatusInvalidParams.Err()
	}
	err := interceptor(context.Background(), "/test", nil, nil, nil, ivoker)
	assert.Error(t, err)
}

func TestSteamClientCircuitBreaker(t *testing.T) {
	interceptor := SteamClientCircuitBreaker()
	assert.NotNil(t, interceptor)

	streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return nil, errcode.StatusInternalServerError.ToRPCErr()
	}
	for i := 0; i < 110; i++ {
		_, err := interceptor(context.Background(), nil, nil, "/test", streamer)
		assert.Error(t, err)
	}

	streamer = func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return nil, errcode.StatusInvalidParams.Err()
	}
	_, err := interceptor(context.Background(), nil, nil, "/test", streamer)
	assert.Error(t, err)
}

func TestUnaryServerCircuitBreaker(t *testing.T) {
	interceptor := UnaryServerCircuitBreaker()
	assert.NotNil(t, interceptor)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errcode.StatusInternalServerError.ToRPCErr()
	}
	for i := 0; i < 110; i++ {
		_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/test"}, handler)
		assert.Error(t, err)
	}

	handler = func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errcode.StatusInvalidParams.Err()
	}
	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/test"}, handler)
	assert.Error(t, err)
}

func TestSteamServerCircuitBreaker(t *testing.T) {
	interceptor := SteamServerCircuitBreaker()
	assert.NotNil(t, interceptor)

	handler := func(srv interface{}, stream grpc.ServerStream) error {
		return errcode.StatusInternalServerError.ToRPCErr()
	}
	for i := 0; i < 110; i++ {
		err := interceptor(nil, nil, &grpc.StreamServerInfo{FullMethod: "/test"}, handler)
		assert.Error(t, err)
	}

	handler = func(srv interface{}, stream grpc.ServerStream) error {
		return errcode.StatusInvalidParams.Err()
	}
	err := interceptor(nil, nil, &grpc.StreamServerInfo{FullMethod: "/test"}, handler)
	assert.Error(t, err)
}
