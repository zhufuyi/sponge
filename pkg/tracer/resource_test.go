package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResource(t *testing.T) {
	resource := NewResource()
	assert.NotNil(t, resource)
}

func TestWithAttributes(t *testing.T) {
	testData := map[string]string{}
	o := new(resourceOptions)
	opt := WithAttributes(testData)
	apply(o, opt)
	assert.Equal(t, testData, o.attributes)
}

func TestWithEnvironment(t *testing.T) {
	testData := "env"
	o := new(resourceOptions)
	opt := WithEnvironment(testData)
	apply(o, opt)
	assert.Equal(t, testData, o.environment)
}

func TestWithServiceName(t *testing.T) {
	testData := "foo"
	o := new(resourceOptions)
	opt := WithServiceName(testData)
	apply(o, opt)
	assert.Equal(t, testData, o.serviceName)
}

func TestWithServiceVersion(t *testing.T) {
	testData := "v1.0"
	o := new(resourceOptions)
	opt := WithServiceVersion(testData)
	apply(o, opt)
	assert.Equal(t, testData, o.serviceVersion)
}

func Test_apply(t *testing.T) {
	testData := "v1.0"
	o := new(resourceOptions)
	opt := WithServiceVersion(testData)
	apply(o, opt)
	assert.Equal(t, testData, o.serviceVersion)
}

func Test_resourceOptionFunc_apply(t *testing.T) {
	testData := "v1.0"
	o := new(resourceOptions)
	opt := WithServiceVersion(testData)
	apply(o, opt)
	assert.Equal(t, testData, o.serviceVersion)
}
