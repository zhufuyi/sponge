package discovery

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/servicerd/registry"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
)

type discovery struct{}

func (d discovery) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	return []*registry.ServiceInstance{}, nil
}

func (d discovery) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return &watcher{}, nil
}

type watcher struct{}

func (w watcher) Next() ([]*registry.ServiceInstance, error) {
	return []*registry.ServiceInstance{}, nil
}

func (w watcher) Stop() error {
	return nil
}

func TestNewBuilder(t *testing.T) {
	b := NewBuilder(&discovery{},
		WithInsecure(false),
		WithTimeout(time.Second),
		DisableDebugLog(),
	)
	assert.NotNil(t, b)
}

func Test_builder_Build(t *testing.T) {
	b := NewBuilder(&discovery{})
	assert.NotNil(t, b)

	u := url.URL{
		Path: "ipv4.single.fake",
	}
	_, err := b.Build(resolver.Target{URL: u}, nil, resolver.BuildOptions{})
	assert.NoError(t, err)
}

func Test_builder_Scheme(t *testing.T) {
	b := NewBuilder(&discovery{})
	assert.NotNil(t, b)
	t.Log(b.Scheme())
}
