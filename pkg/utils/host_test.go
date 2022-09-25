package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAvailablePort(t *testing.T) {
	port, err := GetAvailablePort()
	assert.NoError(t, err)
	t.Log(port)
}

func TestGetHostname(t *testing.T) {
	hostname := GetHostname()
	t.Log(hostname)
}

func TestGetLocalHTTPAddrPairs(t *testing.T) {
	serverAddr, requestAddr := GetLocalHTTPAddrPairs()
	assert.NotEmpty(t, serverAddr)
	assert.NotEmpty(t, requestAddr)
}
