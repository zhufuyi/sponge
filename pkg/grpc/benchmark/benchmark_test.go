package benchmark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	_, err := New("localhost", "test.proto", "Create", nil, 100)
	assert.NoError(t, err)
}

func Test_params_Run(t *testing.T) {
	b, err := New("localhost", "test.proto", "Create", nil, 100)
	assert.NoError(t, err)

	err = b.Run()
	assert.NotNil(t, err)
}
