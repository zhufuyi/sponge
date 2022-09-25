package interceptor

import (
	"testing"

	"github.com/zhufuyi/sponge/pkg/logger"

	"github.com/stretchr/testify/assert"
)

func TestStreamClientLog(t *testing.T) {
	interceptor := StreamClientLog(logger.Get())
	assert.NotNil(t, interceptor)
}

func TestStreamServerCtxTags(t *testing.T) {
	interceptor := StreamServerCtxTags()
	assert.NotNil(t, interceptor)
}

func TestStreamServerLog(t *testing.T) {
	interceptor := StreamServerLog(logger.Get())
	assert.NotNil(t, interceptor)
}

func TestUnaryClientLog(t *testing.T) {
	interceptor := UnaryClientLog(logger.Get())
	assert.NotNil(t, interceptor)
}

func TestUnaryServerCtxTags(t *testing.T) {
	interceptor := UnaryServerCtxTags()
	assert.NotNil(t, interceptor)
}

func TestUnaryServerLog(t *testing.T) {
	interceptor := UnaryServerLog(logger.Get())
	assert.NotNil(t, interceptor)
}

func TestWithLogFields(t *testing.T) {
	testData := map[string]interface{}{"foo": "bar"}
	opt := WithLogFields(testData)
	o := new(logOptions)
	o.apply(opt)
	assert.Equal(t, testData, o.fields)
}

func TestWithLogIgnoreMethods(t *testing.T) {
	testData := "/api.demo.v1"
	opt := WithLogIgnoreMethods(testData)
	o := &logOptions{ignoreMethods: map[string]struct{}{}}
	o.apply(opt)
	assert.Equal(t, struct{}{}, o.ignoreMethods[testData])
}

func Test_defaultLogOptions(t *testing.T) {
	o := defaultLogOptions()
	assert.NotNil(t, o)
}

func Test_logOptions_apply(t *testing.T) {
	testData := map[string]interface{}{"foo": "bar"}
	opt := WithLogFields(testData)
	o := new(logOptions)
	o.apply(opt)
	assert.Equal(t, testData, o.fields)
}
