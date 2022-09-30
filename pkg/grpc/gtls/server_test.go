package gtls

import (
	"testing"

	"github.com/zhufuyi/sponge/pkg/grpc/gtls/certfile"

	"github.com/stretchr/testify/assert"
)

func TestGetServerTLSCredentials(t *testing.T) {
	credentials, err := GetServerTLSCredentials(certfile.Path("one-way/server.crt"), certfile.Path("one-way/server.key"))
	assert.NoError(t, err)
	assert.NotNil(t, credentials)
}

func TestGetServerTLSCredentialsByCA(t *testing.T) {
	credentials, err := GetServerTLSCredentialsByCA(
		certfile.Path("two-way/ca.pem"),
		certfile.Path("two-way/server/server.pem"),
		certfile.Path("two-way/server/server.key"),
	)
	assert.NoError(t, err)
	assert.NotNil(t, credentials)
}
