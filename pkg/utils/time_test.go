package utils

import (
	"testing"
	"time"
)

func TestFormatAndParseDateTime(t *testing.T) {
	layouts := []string{
		DateTimeLayout,
		DateTimeLayoutWithMS,
		DateTimeLayoutWithMSAndTZ,
		TimeLayout,
		DateLayout,
	}

	now := time.Now()

	for _, layout := range layouts {
		str := FormatDateTime(now, layout)
		ti, err := ParseDateTime(str, layout)
		t.Log(str, ti.Second() == now.Second(), err)
	}
}
