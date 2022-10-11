package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewJaegerAgentExporter(t *testing.T) {
	exporter, err := NewJaegerAgentExporter("localhost", "2379")
	assert.NoError(t, err)
	assert.NotNil(t, exporter)
}

func TestNewJaegerExporter(t *testing.T) {
	exporter, err := NewJaegerExporter("http://localhost:14268/api/traces",
		WithUsername("foo"),
		WithPassword("bar"),
	)
	assert.NoError(t, err)
	assert.NotNil(t, exporter)
}

func TestWithPassword(t *testing.T) {
	testData := "123456"
	opt := WithPassword(testData)
	o := new(jaegerOptions)
	o.apply(opt)
	assert.Equal(t, testData, o.password)
}

func TestWithUsername(t *testing.T) {
	testData := "foo"
	opt := WithUsername(testData)
	o := new(jaegerOptions)
	o.apply(opt)
	assert.Equal(t, testData, o.username)
}

func Test_defaultJaegerOptions(t *testing.T) {
	o := defaultJaegerOptions()
	assert.NotNil(t, o)
}

func Test_jaegerOptions_apply(t *testing.T) {
	testData := "foo"
	opt := WithUsername(testData)
	o := new(jaegerOptions)
	o.apply(opt)
	assert.Equal(t, testData, o.username)
}
