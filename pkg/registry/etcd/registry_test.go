package etcd

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zhufuyi/sponge/pkg/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

/*

// 需要连接真实etcd服务测试

import (
	"context"
	"fmt"
	"github.com/zhufuyi/sponge/pkg/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

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

func TestNew(t *testing.T) {
	r := New(&clientv3.Client{},
		WithRegisterTTL(time.Second),
		WithContext(context.Background()),
		WithMaxRetry(3),
		WithNamespace("foo"),
	)
	assert.NotNil(t, r)
}

func TestRegistry_Register(t *testing.T) {
	defer func() {
		recover()
	}()
	r := New(&clientv3.Client{})
	instance := registry.NewServiceInstance("foo", []string{"grpc://127.0.0.1:8282"})
	err := r.Register(context.Background(), instance)
	assert.NoError(t, err)
}

func TestRegistry_Deregister(t *testing.T) {
	defer func() {
		recover()
	}()
	r := New(&clientv3.Client{})
	instance := registry.NewServiceInstance("foo", []string{"grpc://127.0.0.1:8282"})
	err := r.Deregister(context.Background(), instance)
	assert.NoError(t, err)
}
