package grpccli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDial(t *testing.T) {
	_, err := Dial(context.Background(), "localhost:9090")
	assert.NotNil(t, err)
}

func TestDialInsecure(t *testing.T) {
	_, err := DialInsecure(context.Background(), "localhost:9090")
	assert.NoError(t, err)
}

func Test_dial(t *testing.T) {
	_, err := dial(context.Background(), "localhost:9090", true)
	assert.NotNil(t, err)
	_, err = dial(context.Background(), "localhost:9090", false)
	assert.NoError(t, err)
}
