package nacos

import (
	"context"
	"testing"

	"github.com/zhufuyi/sponge/pkg/registry"

	"github.com/stretchr/testify/assert"
)

func newNacosRegistry() *Registry {
	return New(getCli(),
		WithPrefix("/micro"),
		WithWeight(1),
		WithCluster("cluster"),
		WithGroup("dev"),
		WithDefaultKind("grpc"),
	)
}

func TestRegistry(t *testing.T) {
	r := newNacosRegistry()
	instance := registry.NewServiceInstance("foo", []string{"grpc://127.0.0.1:8282"})

	err := r.Register(context.Background(), instance)
	t.Log(err)

	err = r.Deregister(context.Background(), instance)
	t.Log(err)

	_, err = r.GetService(context.Background(), "foo")
	t.Log(err)

	_, err = r.Watch(context.Background(), "foo")
	t.Log(err)
}

func TestRegistry_Register(t *testing.T) {
	r := newNacosRegistry()
	instance := registry.NewServiceInstance("", []string{"grpc://127.0.0.1:8282"})
	err := r.Register(context.Background(), instance)
	assert.Error(t, err)

	instance = registry.NewServiceInstance("foo", []string{"grpc://127.0.0.1:8282"},
		registry.WithMetadata(map[string]string{
			"foo2": "bar2",
		}))
	err = r.Register(context.Background(), instance)
	assert.Error(t, err)

	instance = registry.NewServiceInstance("foo", []string{"127.0.0.1:port"})
	err = r.Register(context.Background(), instance)
	assert.Error(t, err)

	instance = registry.NewServiceInstance("foo", []string{"127.0.0.1"})
	err = r.Register(context.Background(), instance)
	assert.Error(t, err)
}
