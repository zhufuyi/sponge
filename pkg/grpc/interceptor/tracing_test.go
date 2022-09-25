package interceptor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamClientTracing(t *testing.T) {
	interceptor := StreamClientTracing()
	assert.NotNil(t, interceptor)
}

func TestStreamServerTracing(t *testing.T) {
	interceptor := StreamServerTracing()
	assert.NotNil(t, interceptor)
}

func TestUnaryClientTracing(t *testing.T) {
	interceptor := UnaryClientTracing()
	assert.NotNil(t, interceptor)
}

func TestUnaryServerTracing(t *testing.T) {
	interceptor := UnaryServerTracing()
	assert.NotNil(t, interceptor)
}
