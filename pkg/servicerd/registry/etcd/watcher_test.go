package etcd

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/go-dev-frame/sponge/pkg/servicerd/registry"
)

type wt struct{}

func (w wt) Watch(ctx context.Context, key string, opts ...clientv3.OpOption) clientv3.WatchChan {
	c := make(chan clientv3.WatchResponse)
	return c
}

func (w wt) RequestProgress(ctx context.Context) error {
	return nil
}

func (w wt) Close() error {
	return nil
}

func newWatch(first bool) *watcher {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	r := New(&clientv3.Client{})

	return &watcher{
		key:         "foo",
		ctx:         ctx,
		cancel:      cancelFunc,
		watchChan:   make(clientv3.WatchChan),
		watcher:     &wt{},
		kv:          r.kv,
		first:       first,
		serviceName: "host",
	}
}

func Test_watcher_Next(t *testing.T) {
	w := newWatch(false)
	instances, err := w.Next()
	assert.Error(t, err)
	t.Log(instances)

	defer func() { recover() }()
	w = newWatch(true)
	instances, err = w.Next()
	assert.Error(t, err)
	t.Log(instances)
}

func Test_watcher_Stop(t *testing.T) {
	w := newWatch(false)
	err := w.Stop()
	assert.NoError(t, err)
}

func Test_watcher_getInstance(t *testing.T) {
	defer func() { recover() }()

	w := newWatch(false)
	instances, err := w.getInstance()
	assert.NoError(t, err)
	t.Log(instances)
}

func TestService_marshal(t *testing.T) {
	instance := registry.NewServiceInstance("foo", "bar", []string{"grpc://127.0.0.1:8282"})
	v, err := marshal(instance)
	assert.NoError(t, err)

	si, err := unmarshal([]byte(v))
	assert.NoError(t, err)
	assert.Equal(t, instance, si)
}
