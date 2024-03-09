package discovery

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"

	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
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

func Test_discoveryResolver_Close(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
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
		"foo",
		"bar",
		[]string{"grpc://127.0.0.1:8282"},
	)})
	//r.watch()
	//time.Sleep(time.Millisecond * 100)
}

func Test_parseAttributes(t *testing.T) {
	a := parseAttributes(map[string]string{"foo": "bar", "foo2": "bar2"})
	assert.NotNil(t, a)
}

func Test_discoveryResolver_watch(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	r := &discoveryResolver{
		w:                &watcher{},
		cc:               &cliConn{},
		ctx:              ctx,
		cancel:           cancel,
		insecure:         true,
		debugLogDisabled: false,
	}
	defer r.Close()

	r.watch()
	time.Sleep(time.Millisecond * 200)
}

func Test_parseEndpoint(t *testing.T) {
	_, err := parseEndpoint([]string{"grpc://127.0.0.1:8282"}, "grpc", false)
	assert.NoError(t, err)
	_, err = parseEndpoint([]string{"grpc://127.0.0.1:8282"}, "grpc", true)
	assert.NoError(t, err)
	_, err = parseEndpoint(nil, "", true)
	assert.NoError(t, err)
}

func TestIsSecure(t *testing.T) {
	u, err := url.Parse("http://localhost:8080")
	assert.NoError(t, err)

	ok := IsSecure(u)
	assert.Equal(t, false, ok)

	u, _ = url.Parse("http://localhost:8080?isSecure=true")
	ok = IsSecure(u)
	assert.Equal(t, true, ok)
}
