package logger

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAny(t *testing.T) {
	field := Any("key", []int{1, 2, 3})
	assert.NotNil(t, field)
}

func TestBool(t *testing.T) {
	field := Bool("key", true)
	assert.NotNil(t, field)
}

func TestDuration(t *testing.T) {
	field := Duration("key", time.Second)
	assert.NotNil(t, field)
}

func TestErr(t *testing.T) {
	field := Err(errors.New("err"))
	assert.NotNil(t, field)
}

func TestFloat64(t *testing.T) {
	field := Float64("key", 3.14)
	assert.NotNil(t, field)
}

func TestInt(t *testing.T) {
	field := Int("key", 1)
	assert.NotNil(t, field)
}

func TestInt64(t *testing.T) {
	field := Int64("key", 1)
	assert.NotNil(t, field)
}

func TestString(t *testing.T) {
	field := String("key", "bar")
	assert.NotNil(t, field)
}

func TestStringer(t *testing.T) {
	field := Stringer("key", new(st))
	assert.NotNil(t, field)
}

func TestTime(t *testing.T) {
	field := Time("key", time.Now())
	assert.NotNil(t, field)
}

func TestUint(t *testing.T) {
	field := Uint("key", 1)
	assert.NotNil(t, field)
}

func TestUint64(t *testing.T) {
	field := Uint64("key", 1)
	assert.NotNil(t, field)
}

func TestUintptr(t *testing.T) {
	testData := 1
	field := Uintptr("key", uintptr(testData))
	assert.NotNil(t, field)
}

type st struct{}

func (s *st) String() string {
	return "string"
}
