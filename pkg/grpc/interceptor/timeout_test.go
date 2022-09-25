package interceptor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStreamTimeout(t *testing.T) {
	interceptor := StreamTimeout(time.Second)
	assert.NotNil(t, interceptor)
}

func TestUnaryTimeout(t *testing.T) {
	interceptor := UnaryTimeout(time.Second)
	assert.NotNil(t, interceptor)
}

func Test_defaultContextTimeout(t *testing.T) {
	_, cancel := defaultContextTimeout(context.Background())
	if cancel != nil {
		defer cancel()
	}
}
