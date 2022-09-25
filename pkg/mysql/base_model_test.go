package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTableName(t *testing.T) {
	name := GetTableName(&userExample{})
	assert.NotEmpty(t, name)
}
