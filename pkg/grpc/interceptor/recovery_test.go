package interceptor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnaryClientRecovery(t *testing.T) {
	interceptor := UnaryClientRecovery()
	assert.NotNil(t, interceptor)

	err := interceptor(context.Background(), "/test", nil, nil, nil, unaryClientInvoker)
	assert.NoError(t, err)
	err = interceptor(context.Background(), "/test", nil, nil, nil, nil)
	assert.NotNil(t, err)
}

func TestStreamClientRecovery(t *testing.T) {
	interceptor := StreamClientRecovery()
	assert.NotNil(t, interceptor)

	_, err := interceptor(context.Background(), nil, nil, "/test", streamClientFunc)
	assert.NoError(t, err)
	_, err = interceptor(context.Background(), nil, nil, "/test", nil)
	assert.NotNil(t, err)
}

func TestUnaryServerRecovery(t *testing.T) {
	interceptor := UnaryServerRecovery()
	assert.NotNil(t, interceptor)

	_, err := interceptor(context.Background(), nil, unaryServerInfo, unaryServerHandler)
	assert.NoError(t, err)
	_, err = interceptor(context.Background(), nil, unaryServerInfo, nil)
	assert.NotNil(t, err)
}

func TestStreamServerRecovery(t *testing.T) {
	interceptor := StreamServerRecovery()
	assert.NotNil(t, interceptor)

	err := interceptor(nil, newStreamServer(context.Background()), streamServerInfo, streamServerHandler)
	assert.NoError(t, err)
	err = interceptor(nil, newStreamServer(context.Background()), streamServerInfo, nil)
	assert.NotNil(t, err)
}
