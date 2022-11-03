package etcd

import (
	"context"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/servicerd/registry"

	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestNewRegistry(t *testing.T) {
	etcdEndpoints := []string{"127.0.0.1:2379"}
	id := "1"
	instanceName := "serverName"
	instanceEndpoints := []string{"grpc://127.0.0.1:8282"}
	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		iRegistry, serviceInstance, err := NewRegistry(etcdEndpoints, id, instanceName, instanceEndpoints)
		t.Log(err, iRegistry, serviceInstance)
	})
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
	instance := registry.NewServiceInstance("foo", "bar", []string{"grpc://127.0.0.1:8282"})
	err := r.Register(context.Background(), instance)
	assert.NoError(t, err)
}

func TestRegistry_Deregister(t *testing.T) {
	r := newEtcdRegistry()
	instance := registry.NewServiceInstance("foo", "bar", []string{"grpc://127.0.0.1:8282"})
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
