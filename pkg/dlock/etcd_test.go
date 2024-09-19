package dlock

import (
	"fmt"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/etcdcli"
	"go.uber.org/zap"
)

func TestEtcdLock_TryLock(t *testing.T) {
	initLocker := func() Locker {
		return getEtcdLock()
	}
	testLockAndUnlock(initLocker, false, t)
}

func TestEtcdLock_Lock(t *testing.T) {
	initLocker := func() Locker {
		return getEtcdLock()
	}
	testLockAndUnlock(initLocker, true, t)
}

func getEtcdLock() Locker {
	endpoints := []string{"127.0.0.1:2379"}
	cli, err := etcdcli.Init(endpoints,
		etcdcli.WithDialTimeout(time.Second*2),
		etcdcli.WithAuth("", ""),
		etcdcli.WithAutoSyncInterval(0),
		etcdcli.WithLog(zap.NewNop()),
	)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	locker, err := NewEtcd(cli, "sponge/dlock", 10)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return locker
}
