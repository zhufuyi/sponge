package interceptor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamServerRecovery(t *testing.T) {
	interceptor := StreamServerRecovery()
	assert.NotNil(t, interceptor)
}

func TestUnaryServerRecovery(t *testing.T) {
	interceptor := UnaryServerRecovery()
	assert.NotNil(t, interceptor)
}
