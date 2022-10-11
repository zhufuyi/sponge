package gofile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindSubBytes(t *testing.T) {
	testData := []byte(`start1234567890end`)
	val := FindSubBytes(testData, []byte("start"), []byte("end"))
	assert.Equal(t, testData, val)

	val = FindSubBytes(testData, []byte("end"), []byte("start"))
	assert.Empty(t, val)
}

func TestFindSubBytesNotIn(t *testing.T) {
	testData := []byte(`start1234567890end`)
	val := FindSubBytesNotIn(testData, []byte("start"), []byte("end"))
	assert.Equal(t, []byte("1234567890"), val)

	val = FindSubBytesNotIn(testData, []byte("end"), []byte("start"))
	assert.Empty(t, val)
}
