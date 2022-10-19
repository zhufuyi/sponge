package nacos

import (
	"context"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/stretchr/testify/assert"
	"github.com/zhufuyi/sponge/pkg/nacoscli"
	"testing"
	"time"
)

func getCli() naming_client.INamingClient {
	params := &nacoscli.Params{
		IPAddr:      "127.0.0.1",
		Port:        8448,
		NamespaceID: "de7b176e-91cd-49a3-ac83-beb725979775",
		Group:       "dev",
		DataID:      "user-srv.yml",
		Format:      "yaml",
	}
	namingClient, err := nacoscli.NewNamingClient(params)
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
