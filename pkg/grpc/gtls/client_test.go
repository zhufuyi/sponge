package gtls

import (
	"testing"

	"github.com/zhufuyi/sponge/pkg/grpc/gtls/certfile"

	"github.com/stretchr/testify/assert"
)

func TestGetClientTLSCredentials(t *testing.T) {
	credentials, err := GetClientTLSCredentials("localhost", certfile.Path("one-way/server.crt"))
	assert.NoError(t, err)
	assert.NotNil(t, credentials)

	_, err = GetClientTLSCredentials("localhost", certfile.Path("one-way/notfound.crt"))
	assert.Error(t, err)
}

func TestGetClientTLSCredentialsByCA(t *testing.T) {
	credentials, err := GetClientTLSCredentialsByCA(
		"localhost",
		certfile.Path("two-way/ca.pem"),
		certfile.Path("two-way/client/client.pem"),
		certfile.Path("two-way/client/client.key"),
	)
	assert.NoError(t, err)
	assert.NotNil(t, credentials)

	_, err = GetClientTLSCredentialsByCA(
		"localhost",
		certfile.Path("two-way/ca.pem"),
		certfile.Path("two-way/client/notfound.pem"),
		certfile.Path("two-way/client/notfound.key"),
	)
	assert.Error(t, err)

	_, err = GetClientTLSCredentialsByCA(
		"localhost",
		certfile.Path("two-way/notfound.pem"),
		certfile.Path("two-way/client/client.pem"),
		certfile.Path("two-way/client/client.key"),
	)
	assert.Error(t, err)
}
