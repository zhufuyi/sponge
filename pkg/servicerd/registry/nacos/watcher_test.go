package nacos

import (
	"context"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/nacoscli"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/stretchr/testify/assert"
)

func getCli() naming_client.INamingClient {
	var (
		ipAddr      = "192.168.3.37"
		port        = 8848
		namespaceID = "3454d2b5-2455-4d0e-bf6d-e033b086bb4c"
	)
	namingClient, err := nacoscli.NewNamingClient(ipAddr, port, namespaceID)
	if err != nil {
		panic(err)
	}

	return namingClient
}

func newWatch() *watcher {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	wt := &watcher{
		serviceName: "host",
		clusters:    []string{"bar"},
		groupName:   "foo",
		ctx:         ctx,
		cancel:      cancelFunc,
		watchChan:   make(chan struct{}),
		cli:         getCli(),
		kind:        "host",
	}

	return wt
}

func Test_newWatcher(t *testing.T) {
	defer func() { recover() }()
	_, _ = newWatcher(context.Background(), getCli(), "host", "host", "foo", []string{"bar"})
}

func Test_watcher(t *testing.T) {
	defer func() { recover() }()
	_, _ = newWatcher(context.Background(), getCli(), "host", "host", "foo", []string{"bar"})

	w := newWatch()
	_, err := w.Next()
	t.Log(err)

	err = w.Stop()
	assert.NoError(t, err)
}
