package interceptor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnaryClientHystrix(t *testing.T) {
	interceptor := UnaryClientHystrix("demo")
	assert.NotNil(t, interceptor)
}

func TestSteamClientHystrix(t *testing.T) {
	interceptor := SteamClientHystrix("demo")
	assert.NotNil(t, interceptor)
}
