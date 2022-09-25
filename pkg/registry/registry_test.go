package registry

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewServiceInstance(t *testing.T) {
	s := NewServiceInstance("demo", []string{"grpc://127.0.0.1:9090"},
		WithID("1"),
		WithVersion("v1.0.0"),
		WithMetadata(map[string]string{"foo": "bar"}),
	)
	assert.NotNil(t, s)
}
