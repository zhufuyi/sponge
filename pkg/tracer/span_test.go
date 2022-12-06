package tracer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSpan(t *testing.T) {
	SetTraceName("foo")

	tags := map[string]interface{}{
		"foo1":  nil,
		"foo2":  true,
		"foo3":  "bar",
		"foo4":  []string{"bar"},
		"foo5":  1,
		"foo6":  []int{1},
		"foo7":  int64(1),
		"foo8":  []int64{1},
		"foo9":  3.14,
		"foo10": []float64{3.14},
		"foo11": map[string]string{"foo": "bar"},
	}
	_, span := NewSpan(context.Background(), "fooSpan", tags)
	assert.NotNil(t, span)
}
