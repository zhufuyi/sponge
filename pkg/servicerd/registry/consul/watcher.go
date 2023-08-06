package consul

import (
	"context"

	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
)

type watcher struct {
	event chan struct{}
	set   *serviceSet

	// for cancel
	ctx    context.Context
	cancel context.CancelFunc
}

func (w *watcher) Next() (services []*registry.ServiceInstance, err error) {
	select {
	case <-w.ctx.Done():
		err = w.ctx.Err()
	case <-w.event:
	}

	ss, ok := w.set.services.Load().([]*registry.ServiceInstance)

	if ok {
		services = append(services, ss...)
	}
	return //nolint
}

func (w *watcher) Stop() error {
	w.cancel()
	w.set.lock.Lock()
	defer w.set.lock.Unlock()
	delete(w.set.watcher, w)
	return nil
}
