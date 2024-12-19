package consul

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-dev-frame/sponge/pkg/servicerd/registry"
)

func newServiceSet() *serviceSet {
	return &serviceSet{
		serviceName: "foo",
		watcher:     map[*watcher]struct{}{},
		services:    &atomic.Value{},
		lock:        sync.RWMutex{},
	}
}

func TestServiceSet_broadcast(t *testing.T) {
	ss := newServiceSet()
	instance := registry.NewServiceInstance("foo", "bar", []string{"grpc://127.0.0.1:8282"})
	ss.broadcast([]*registry.ServiceInstance{instance})
}

func newWatch() *watcher {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	wt := &watcher{
		event:  make(chan struct{}),
		set:    newServiceSet(),
		ctx:    ctx,
		cancel: cancelFunc,
	}

	return wt
}

func Test_watcher(t *testing.T) {
	w := newWatch()

	_, err := w.Next()
	t.Log(err)

	err = w.Stop()
	t.Log(err)
}
