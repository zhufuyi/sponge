package consulcli

import (
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	addr := "192.168.3.37:8500"
	cli, err := Init(addr,
		WithScheme("http"),
		WithWaitTime(time.Second*2),
		WithDatacenter(""),
	)
	t.Log(err, cli)
}
