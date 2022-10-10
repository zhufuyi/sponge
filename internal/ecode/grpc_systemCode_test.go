package ecode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAny(t *testing.T) {
	detail := Any("foo", "bar")
	assert.Contains(t, detail.String(), "foo")

	detail = Any("foo1", "bar1")
	assert.NotContains(t, detail.String(), "any-key")
}
