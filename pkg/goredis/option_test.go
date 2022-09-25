package goredis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithEnableTrace(t *testing.T) {
	opt := WithEnableTrace()
	o := new(options)
	o.apply(opt)
	assert.Equal(t, true, o.enableTrace)
}

func Test_defaultOptions(t *testing.T) {
	o := defaultOptions()
	assert.NotNil(t, o)
}
