package resolve

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

var r = &ResolverBuilder{
	scheme:      "grpc",
	serviceName: "demo",
	addrs:       []string{"localhost:8282"},
}

func TestRegister(t *testing.T) {
	s := Register(r.scheme, r.serviceName, r.addrs)
	assert.Equal(t, true, strings.Contains(s, r.serviceName))
}

func TestResolverBuilder_Build(t *testing.T) {
	c := &clientConn{}
	_, err := r.Build(resolver.Target{}, c, resolver.BuildOptions{})
	assert.NoError(t, err)
}

func TestResolverBuilder_Scheme(t *testing.T) {
	str := r.Scheme()
	assert.NotEmpty(t, str)
}

func Test_blResolver_Close(t *testing.T) {
	c := &clientConn{}
	b, err := r.Build(resolver.Target{}, c, resolver.BuildOptions{})
	assert.NoError(t, err)

	b.Close()
}

func Test_blResolver_ResolveNow(t *testing.T) {
	c := &clientConn{}
	b, err := r.Build(resolver.Target{}, c, resolver.BuildOptions{})
	assert.NoError(t, err)

	b.ResolveNow(struct{}{})
}

func Test_blResolver_start(t *testing.T) {
	b := &blResolver{
		target:     resolver.Target{},
		cc:         &clientConn{},
		addrsStore: make(map[string][]string),
	}
	b.start()
}

type clientConn struct{}

func (c clientConn) UpdateState(state resolver.State) error  { return nil }
func (c clientConn) ReportError(err error)                   {}
func (c clientConn) NewAddress(addresses []resolver.Address) {}
func (c clientConn) NewServiceConfig(serviceConfig string)   {}
func (c clientConn) ParseServiceConfig(serviceConfigJSON string) *serviceconfig.ParseResult {
	return &serviceconfig.ParseResult{}
}
