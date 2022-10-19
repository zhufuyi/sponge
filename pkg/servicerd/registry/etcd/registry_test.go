package etcd

import (
	"context"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/servicerd/registry"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

/*
// need real etcd service test
func TestGRPCSeverRegistry(t *testing.T) {
	instance := registry.NewServiceInstance("foo", []string{"grpc://127.0.0.1:8282"})

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.3.37:2379"},
		DialTimeout: 3 * time.Second,
		DialOptions: []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	etcdRegistry := New(cli)
	ctx := context.Background()

	err = etcdRegistry.Register(ctx, instance)
	if err != nil {
		t.Fatal(err)
	}

	instances, err := etcdRegistry.GetService(ctx, instance.Name)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("register %+v", instances[0])

	time.Sleep(3 * time.Second)

	t.Log("deregister")
	err = etcdRegistry.Deregister(ctx, instance)
	if err != nil {
		t.Fatal(err)
	}

}

func TestRegistry(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.3.37:2379"},
		DialTimeout: time.Second, DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()
	s := &registry.ServiceInstance{
		ID:   "0",
		Name: "helloworld",
	}

	r := New(client)
	w, err := r.Watch(ctx, s.Name)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = w.Stop()
	}()
	go func() {
		for {
			res, err1 := w.Next()
			if err1 != nil {
				return
			}
			t.Logf("watch: %d", len(res))
			for _, r := range res {
				t.Logf("next: %+v", r)
			}
		}
	}()
	time.Sleep(time.Second)

	if err1 := r.Register(ctx, s); err1 != nil {
		t.Fatal(err1)
	}
	time.Sleep(time.Second)

	res, err := r.GetService(ctx, s.Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 && res[0].Name != s.Name {
		t.Errorf("not expected: %+v", res)
	}

	if err1 := r.Deregister(ctx, s); err1 != nil {
		t.Fatal(err1)
	}
	time.Sleep(time.Second)

	res, err = r.GetService(ctx, s.Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 0 {
		t.Errorf("not expected empty")
	}
}

func TestHeartBeat(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.3.37:2379"},
		DialTimeout: time.Second, DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()
	s := &registry.ServiceInstance{
		ID:   "0",
		Name: "helloworld",
	}

	go func() {
		r := New(client)
		w, err1 := r.Watch(ctx, s.Name)
		if err1 != nil {
			return
		}
		defer func() {
			_ = w.Stop()
		}()
		for {
			res, err2 := w.Next()
			if err2 != nil {
				return
			}
			t.Logf("watch: %d", len(res))
			for _, r := range res {
				t.Logf("next: %+v", r)
			}
		}
	}()
	time.Sleep(time.Second)

	// new a server
	r := New(client,
		WithRegisterTTL(2*time.Second),
		WithMaxRetry(5),
	)

	key := fmt.Sprintf("%s/%s/%s", r.opts.namespace, s.Name, s.ID)
	value, _ := marshal(s)
	r.lease = clientv3.NewLease(r.client)
	leaseID, err := r.registerWithKV(ctx, key, value)
	if err != nil {
		t.Fatal(err)
	}

	// wait for lease expired
	time.Sleep(3 * time.Second)

	res, err := r.GetService(ctx, s.Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 0 {
		t.Errorf("not expected empty")
	}

	go r.heartBeat(ctx, leaseID, key, value)

	time.Sleep(time.Second)
	res, err = r.GetService(ctx, s.Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) == 0 {
		t.Errorf("reconnect failed")
	}
}
*/

func TestNewRegistry(t *testing.T) {
	etcdEndpoints := []string{"127.0.0.1:2379"}
	instanceName := "serverName"
	instanceEndpoints := []string{"grpc://127.0.0.1:8282"}
	iRegistry, serviceInstance, err := NewRegistry(etcdEndpoints, instanceName, instanceEndpoints)
	t.Log(err, iRegistry, serviceInstance)
}

type lease struct{}

func (l lease) Grant(ctx context.Context, ttl int64) (*clientv3.LeaseGrantResponse, error) {
	return &clientv3.LeaseGrantResponse{}, nil
}

func (l lease) Revoke(ctx context.Context, id clientv3.LeaseID) (*clientv3.LeaseRevokeResponse, error) {
	return &clientv3.LeaseRevokeResponse{}, nil
}

func (l lease) TimeToLive(ctx context.Context, id clientv3.LeaseID, opts ...clientv3.LeaseOption) (*clientv3.LeaseTimeToLiveResponse, error) {
	return &clientv3.LeaseTimeToLiveResponse{}, nil
}

func (l lease) Leases(ctx context.Context) (*clientv3.LeaseLeasesResponse, error) {
	return &clientv3.LeaseLeasesResponse{}, nil
}

func (l lease) KeepAlive(ctx context.Context, id clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	c := make(chan *clientv3.LeaseKeepAliveResponse)
	return c, nil
}

func (l lease) KeepAliveOnce(ctx context.Context, id clientv3.LeaseID) (*clientv3.LeaseKeepAliveResponse, error) {
	return &clientv3.LeaseKeepAliveResponse{}, nil
}

func (l lease) Close() error {
	return nil
}

type kv struct{}

func (k kv) Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	return &clientv3.PutResponse{}, nil
}

func (k kv) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	return &clientv3.GetResponse{}, nil
}

func (k kv) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	return &clientv3.DeleteResponse{}, nil
}

func (k kv) Compact(ctx context.Context, rev int64, opts ...clientv3.CompactOption) (*clientv3.CompactResponse, error) {
	return &clientv3.CompactResponse{}, nil
}

func (k kv) Do(ctx context.Context, op clientv3.Op) (clientv3.OpResponse, error) {
	return clientv3.OpResponse{}, nil
}

func (k kv) Txn(ctx context.Context) clientv3.Txn {
	return nil
}

func newEtcdRegistry() *Registry {
	r := New(&clientv3.Client{Lease: &lease{}, KV: &kv{}},
		WithRegisterTTL(time.Second),
		WithContext(context.Background()),
		WithMaxRetry(3),
		WithNamespace("foo"),
	)
	r.lease = &lease{}
	r.kv = &kv{}
	return r
}

func TestRegistry_Register(t *testing.T) {
	defer func() { recover() }()

	r := newEtcdRegistry()
	instance := registry.NewServiceInstance("foo", []string{"grpc://127.0.0.1:8282"})
	err := r.Register(context.Background(), instance)
	assert.NoError(t, err)
}

func TestRegistry_Deregister(t *testing.T) {
	r := newEtcdRegistry()
	instance := registry.NewServiceInstance("foo", []string{"grpc://127.0.0.1:8282"})
	err := r.Deregister(context.Background(), instance)
	assert.NoError(t, err)
}

func TestRegistry_GetService(t *testing.T) {
	r := newEtcdRegistry()
	_, err := r.GetService(context.Background(), "foo")
	assert.NoError(t, err)
}

func TestRegistry_registerWithKV(t *testing.T) {
	r := newEtcdRegistry()
	_, err := r.registerWithKV(context.Background(), "foo", "bar")
	assert.NoError(t, err)
}

func TestRegistry_heartBeat(t *testing.T) {
	r := newEtcdRegistry()
	go r.heartBeat(context.Background(), 1, "foo", "bar")
	time.Sleep(time.Second)
}

func TestRegistry_retry(t *testing.T) {
	r := newEtcdRegistry()
	ctx := context.Background()
	leaseID := clientv3.LeaseID(0)
	kac, _ := r.client.KeepAlive(ctx, leaseID)
	go r.retry(ctx, leaseID, "foo", "bar", kac)
	time.Sleep(time.Second)
}
