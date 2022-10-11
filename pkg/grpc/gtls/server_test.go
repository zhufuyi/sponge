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

	_, err = GetServerTLSCredentials(certfile.Path("one-way/notfound.crt"), certfile.Path("one-way/notfound.key"))
	assert.Error(t, err)
}

func TestGetServerTLSCredentialsByCA(t *testing.T) {
	credentials, err := GetServerTLSCredentialsByCA(
		certfile.Path("two-way/ca.pem"),
		certfile.Path("two-way/server/server.pem"),
		certfile.Path("two-way/server/server.key"),
	)
	assert.NoError(t, err)
	assert.NotNil(t, credentials)

	_, err = GetServerTLSCredentialsByCA(
		certfile.Path("two-way/ca.pem"),
		certfile.Path("two-way/server/notfound.pem"),
		certfile.Path("two-way/server/notfound.key"),
	)
	assert.Error(t, err)

	_, err = GetServerTLSCredentialsByCA(
		certfile.Path("two-way/notfound.pem"),
		certfile.Path("two-way/server/server.pem"),
		certfile.Path("two-way/server/server.key"),
	)
	assert.Error(t, err)
}
