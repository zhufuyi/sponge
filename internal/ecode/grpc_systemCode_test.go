package ecode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAny(t *testing.T) {
	detail := Any("foo", "bar")
	assert.Equal(t, "foo: {bar}", detail.String())
}
