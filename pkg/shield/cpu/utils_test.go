package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	_, err := readFile("notfound")
	assert.Error(t, err)
	_, err = readFile("utils.go")
	assert.NoError(t, err)

	_, err = parseUint("-100")
	assert.Nil(t, err)
	_, err = parseUint("abc")
	assert.Error(t, err)
	_, err = parseUint("100")

	_, err = ParseUintList("")
	assert.Nil(t, err)

	_, err = readLines("notfound")
	assert.Error(t, err)
	_, err = readLinesOffsetN("test.data", 0, -1)
	assert.NoError(t, err)
}

func TestParseUintList(t *testing.T) {
	testData := []string{
		"",
		"7",
		"1-6",
		"0,3-4,7,8-10",
		"0-0,0,1-7",
		"03,1-3",
		"3,2,1",
		"0-2,3,1",
	}

	for _, val := range testData {
		_, err := ParseUintList(val)
		assert.NoError(t, err)
	}
}
