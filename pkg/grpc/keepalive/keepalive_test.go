package keepalive

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClientKeepAlive(t *testing.T) {
	alive := ClientKeepAlive()
	assert.NotNil(t, alive)
}

func TestServerKeepAlive(t *testing.T) {
	alives := ServerKeepAlive()
	assert.Equal(t, 2, len(alives))
}
