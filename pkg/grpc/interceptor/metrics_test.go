package interceptor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamClientMetrics(t *testing.T) {
	interceptor := StreamClientMetrics()
	assert.NotNil(t, interceptor)
}

func TestStreamServerMetrics(t *testing.T) {
	interceptor := StreamServerMetrics()
	assert.NotNil(t, interceptor)
}

func TestUnaryClientMetrics(t *testing.T) {
	interceptor := UnaryClientMetrics()
	assert.NotNil(t, interceptor)
}

func TestUnaryServerMetrics(t *testing.T) {
	interceptor := UnaryServerMetrics()
	assert.NotNil(t, interceptor)
}
