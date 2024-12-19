package consul

import (
	"context"
	"testing"

	"github.com/hashicorp/consul/api"

	"github.com/go-dev-frame/sponge/pkg/servicerd/registry"
)

func getConsulClient() *Client {
	consulClient, err := api.NewClient(&api.Config{})
	if err != nil {
		panic(err)
	}

	return NewClient(consulClient)
}

func TestConsulClient(t *testing.T) {
	cli := getConsulClient()

	instance := registry.NewServiceInstance("foo", "bar", []string{"grpc://127.0.0.1:8282"})
	err := cli.Register(context.Background(), instance, false)
	t.Log(err)

	_, _, err = cli.Service(context.Background(), "foo", 1, false)
	t.Log(err)

	err = cli.Deregister(context.Background(), "1")
	t.Log(err)
}
