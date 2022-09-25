package discovery

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zhufuyi/sponge/pkg/registry"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
	"net/url"
	"testing"
	"time"
)

type cliConn struct {
}

func (c cliConn) UpdateState(state resolver.State) error {
	return nil
}

func (c cliConn) ReportError(err error) {}

func (c cliConn) NewAddress(addresses []resolver.Address) {}

func (c cliConn) NewServiceConfig(serviceConfig string) {}

func (c cliConn) ParseServiceConfig(serviceConfigJSON string) *serviceconfig.ParseResult {
	return &serviceconfig.ParseResult{}
}

func TestIsSecure(t *testing.T) {
	u, err := url.Parse("http://localhost:8080")
	assert.NoError(t, err)

	ok := IsSecure(u)
	assert.Equal(t, false, ok)
}

func Test_discoveryResolver_Close(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	r := &discoveryResolver{
		w:                &watcher{},
		cc:               &cliConn{},
		ctx:              ctx,
		cancel:           cancel,
		insecure:         true,
		debugLogDisabled: false,
	}
	defer r.Close()

	r.ResolveNow(resolver.ResolveNowOptions{})
	r.update([]*registry.ServiceInstance{registry.NewServiceInstance(
		"demo",
		[]string{"grpc://127.0.0.1:9090"},
	)})
	r.watch()
}

func Test_parseAttributes(t *testing.T) {
	a := parseAttributes(map[string]string{"foo": "bar"})
	assert.NotNil(t, a)
}
