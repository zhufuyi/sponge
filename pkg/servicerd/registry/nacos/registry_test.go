package nacos

import (
	"context"
	"testing"

	"github.com/zhufuyi/sponge/pkg/servicerd/registry"

	"github.com/stretchr/testify/assert"
)

func TestNewRegistry(t *testing.T) {
	nacosIPAddr := "192.168.3.37"
	nacosPort := 8848
	nacosNamespaceID := "3454d2b5-2455-4d0e-bf6d-e033b086bb4c"

	id := "serverName_192.168.3.37"
	instanceName := "serverName"
	instanceEndpoints := []string{"grpc://192.168.3.27:8282"}

	iRegistry, instance, err := NewRegistry(nacosIPAddr, nacosPort, nacosNamespaceID, id, instanceName, instanceEndpoints)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(iRegistry, instance)
}

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
	instance := registry.NewServiceInstance("foo", "bar", []string{"grpc://127.0.0.1:8282"})

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
	instance := registry.NewServiceInstance("", "", []string{"grpc://127.0.0.1:8282"})
	err := r.Register(context.Background(), instance)
	assert.Error(t, err)

	instance = registry.NewServiceInstance("foo", "bar", []string{"grpc://127.0.0.1:8282"},
		registry.WithMetadata(map[string]string{
			"foo2": "bar2",
		}))
	err = r.Register(context.Background(), instance)
	assert.NoError(t, err)

	instance = registry.NewServiceInstance("foo", "bar", []string{"127.0.0.1:port"})
	err = r.Register(context.Background(), instance)
	assert.Error(t, err)

	instance = registry.NewServiceInstance("foo", "bar", []string{"127.0.0.1"})
	err = r.Register(context.Background(), instance)
	assert.Error(t, err)
}
