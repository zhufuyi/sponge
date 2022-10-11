package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntToStr(t *testing.T) {
	val := IntToStr(1)
	assert.Equal(t, "1", val)
}

func TestStrToFloat32(t *testing.T) {
	val := StrToFloat32("1")
	assert.Equal(t, float32(1), val)
}

func TestStrToFloat32E(t *testing.T) {
	val, err := StrToFloat32E("1")
	assert.NoError(t, err)
	assert.Equal(t, float32(1), val)
}

func TestStrToFloat64(t *testing.T) {
	val := StrToFloat64("1")
	assert.Equal(t, 1.0, val)
}

func TestStrToFloat64E(t *testing.T) {
	val, err := StrToFloat64E("1")
	assert.NoError(t, err)
	assert.Equal(t, 1.0, val)
}

func TestStrToInt(t *testing.T) {
	val := StrToInt("1")
	assert.Equal(t, 1, val)
}

func TestStrToIntE(t *testing.T) {
	val, err := StrToIntE("1")
	assert.NoError(t, err)
	assert.Equal(t, 1, val)
}

func TestStrToUint32(t *testing.T) {
	val := StrToUint32("1")
	assert.Equal(t, uint32(1), val)
}

func TestStrToUint32E(t *testing.T) {
	val, err := StrToUint32E("1")
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), val)
}

func TestStrToUint64(t *testing.T) {
	val := StrToUint64("1")
	assert.Equal(t, uint64(1), val)
}

func TestStrToUint64E(t *testing.T) {
	val, err := StrToUint64E("1")
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), val)
}

func TestUint64ToStr(t *testing.T) {
	val := Uint64ToStr(1)
	assert.Equal(t, "1", val)
}

func TestInt64ToStr(t *testing.T) {
	val := Int64ToStr(1)
	assert.Equal(t, "1", val)
}
