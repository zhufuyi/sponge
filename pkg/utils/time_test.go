package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatAndParseDateTime(t *testing.T) {
	now := time.Now()

	s := FormatDateTimeLayout(now)
	ti, err := ParseDateTimeLayout(s)
	assert.NoError(t, err)
	assert.Equal(t, ti.Second() == now.Second(), true)

	s = FormatDateTimeLayoutWithMS(now)
	ti, err = ParseDateTimeLayoutWithMS(s)
	assert.NoError(t, err)
	assert.Equal(t, ti.Second() == now.Second(), true)

	s = FormatDateTimeRFC3339(now)
	ti, err = ParseDateTimeRFC3339(s)
	assert.NoError(t, err)
	assert.Equal(t, ti.Second() == now.Second(), true)

	s = FormatDateTimeLayoutWithMSAndTZ(now)
	ti, err = ParseDateTimeLayoutWithMSAndTZ(s)
	assert.NoError(t, err)
	assert.Equal(t, ti.Second() == now.Second(), true)

	ti, err = ParseDateTimeLayout(TimeLayout)
	assert.Error(t, err)
	ti, err = ParseDateTimeLayout(DateLayout)
	assert.Error(t, err)
}
