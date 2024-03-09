package keepalive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientKeepAlive(t *testing.T) {
	alive := ClientKeepAlive()
	assert.NotNil(t, alive)
}

func TestServerKeepAlive(t *testing.T) {
	alives := ServerKeepAlive()
	assert.Equal(t, 2, len(alives))
}
