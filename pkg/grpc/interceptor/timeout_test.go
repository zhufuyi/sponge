package interceptor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnaryTimeout(t *testing.T) {
	interceptor := UnaryTimeout(time.Second)
	assert.NotNil(t, interceptor)

	err := interceptor(context.Background(), "/test", nil, nil, nil, unaryClientInvoker)
	assert.NoError(t, err)
}

func TestStreamTimeout(t *testing.T) {
	interceptor := StreamTimeout(time.Second)
	assert.NotNil(t, interceptor)

	_, err := interceptor(context.Background(), nil, nil, "/test", streamClientFunc)
	assert.NoError(t, err)
}

func Test_defaultContextTimeout(t *testing.T) {
	_, cancel := defaultContextTimeout(context.Background())
	if cancel != nil {
		defer cancel()
	}
}
