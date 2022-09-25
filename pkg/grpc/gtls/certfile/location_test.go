package certfile

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPath(t *testing.T) {
	testData := "README.md"
	file := Path(testData)
	assert.Equal(t, true, strings.Contains(file, testData))
}
