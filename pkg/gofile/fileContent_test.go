package gofile

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindSubBytes(t *testing.T) {
	testData := []byte(`start1234567890end`)
	val := FindSubBytes(testData, []byte("start"), []byte("end"))
	assert.Equal(t, testData, val)
}

func TestFindSubBytesNotIn(t *testing.T) {
	testData := []byte(`start1234567890end`)
	val := FindSubBytesNotIn(testData, []byte("start"), []byte("end"))
	assert.Equal(t, []byte("1234567890"), val)
}
