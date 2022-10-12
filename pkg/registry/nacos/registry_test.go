package nacos

import (
	"context"
	"github.com/zhufuyi/sponge/pkg/registry"
	"testing"
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
