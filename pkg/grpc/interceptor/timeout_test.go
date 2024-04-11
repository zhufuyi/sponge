package interceptor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnaryClientTimeout(t *testing.T) {
	interceptor := UnaryClientTimeout(time.Millisecond)
	assert.NotNil(t, interceptor)
	interceptor = UnaryClientTimeout(time.Second)
	assert.NotNil(t, interceptor)

	err := interceptor(context.Background(), "/test", nil, nil, nil, unaryClientInvoker)
	assert.NoError(t, err)
}

func TestStreamClientTimeout(t *testing.T) {
	interceptor := StreamClientTimeout(time.Millisecond)
	assert.NotNil(t, interceptor)
	interceptor = StreamClientTimeout(time.Second)
	assert.NotNil(t, interceptor)

	_, err := interceptor(context.Background(), nil, nil, "/test", streamClientFunc)
	assert.NoError(t, err)
}
