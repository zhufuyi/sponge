package consulcli

import (
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
)

func TestInit(t *testing.T) {
	addr := "192.168.3.37:8500"
	cli, err := Init(addr,
		WithScheme("http"),
		WithWaitTime(time.Second*2),
		WithDatacenter(""),
	)
	t.Log(err, cli)

	cli, err = Init("", WithConfig(&api.Config{
		Address:    addr,
		Scheme:     "http",
		WaitTime:   time.Second * 2,
		Datacenter: "",
	}))
	t.Log(err, cli)
}
