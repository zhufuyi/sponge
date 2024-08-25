package interceptor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/zhufuyi/sponge/pkg/container/group"
	"github.com/zhufuyi/sponge/pkg/errcode"
	"github.com/zhufuyi/sponge/pkg/shield/circuitbreaker"
	"google.golang.org/grpc/codes"
)

func TestUnaryClientCircuitBreaker(t *testing.T) {
	interceptor := UnaryClientCircuitBreaker(
		WithGroup(group.NewGroup(func() interface{} {
			return circuitbreaker.NewBreaker()
		})),
		WithValidCode(codes.PermissionDenied),
	)

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
	interceptor := StreamClientCircuitBreaker()
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
	degradeHandler := func(ctx context.Context, req interface{}) (reply interface{}, error error) {
		return "degrade", errcode.StatusSuccess.ToRPCErr()
	}
	interceptor := UnaryServerCircuitBreaker(WithUnaryServerDegradeHandler(degradeHandler))
	assert.NotNil(t, interceptor)

	count := 0
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		count++
		if count%2 == 0 {
			return nil, errcode.StatusSuccess.ToRPCErr()
		}
		return nil, errcode.StatusInternalServerError.ToRPCErr()
	}

	successCount, failCount, degradeCount := 0, 0, 0
	for i := 0; i < 1000; i++ {
		reply, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/test"}, handler)
		if err != nil {
			failCount++
			continue
		}
		if reply == "degrade" {
			degradeCount++
		} else {
			successCount++
		}
	}
	t.Logf("successCount: %d, failCount: %d, degradeCount: %d", successCount, failCount, degradeCount)

	handler = func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errcode.StatusInvalidParams.Err()
	}
	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/test"}, handler)
	t.Log(err)
}

func TestSteamServerCircuitBreaker(t *testing.T) {
	interceptor := StreamServerCircuitBreaker()
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
