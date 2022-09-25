package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	err := Init("empty file")
	assert.Error(t, err)

	c := Get()
	assert.NotNil(t, c)

	str := Show()
	assert.NotEmpty(t, str)
}

func TestPath(t *testing.T) {
	path := Path("conf.yml")
	t.Log(path)
}
