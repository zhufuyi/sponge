package tracer

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	exporter, err := newExporter(os.Stdout)
	assert.NoError(t, err)
	resource := NewResource()
	Init(exporter, resource)
}

func TestClose(t *testing.T) {
	exporter, err := newExporter(os.Stdout)
	assert.NoError(t, err)
	resource := NewResource()
	Init(exporter, resource)
	_ = Close(context.Background())
}
