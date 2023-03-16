package interceptor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestUnaryServerRateLimit(t *testing.T) {
	interceptor := UnaryServerRateLimit(
		WithWindow(time.Second*10),
		WithBucket(200),
		WithCPUThreshold(500),
		WithCPUQuota(0.5),
	)
	assert.NotNil(t, interceptor)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}
	_, err := interceptor(nil, nil, nil, handler)
	assert.NoError(t, err)
}

func TestStreamServerRateLimit(t *testing.T) {
	interceptor := StreamServerRateLimit()
	assert.NotNil(t, interceptor)

	handler := func(srv interface{}, stream grpc.ServerStream) error {
		return nil
	}
	err := interceptor(nil, nil, nil, handler)
	assert.NoError(t, err)
}
