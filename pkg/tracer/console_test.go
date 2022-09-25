package tracer

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewConsoleExporter(t *testing.T) {
	exporter, err := NewConsoleExporter()
	assert.NoError(t, err)
	assert.NotNil(t, exporter)
}

func TestNewFileExporter(t *testing.T) {
	exporter, file, err := NewFileExporter("demo")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, exporter)
	_ = file.Close()
	_ = os.RemoveAll("demo")
}

func Test_newExporter(t *testing.T) {
	exporter, err := newExporter(os.Stdout)
	assert.NoError(t, err)
	assert.NotNil(t, exporter)
}
