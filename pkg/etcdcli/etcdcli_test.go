package etcdcli

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	endpoints := []string{"192.168.3.37:2379"}
	cli, err := Init(endpoints,
		WithDialTimeout(time.Second*2),
		WithAuth("", ""),
		WithAutoSyncInterval(0),
		WithLog(zap.NewNop()),
	)
	t.Log(err, cli)

	// test error
	_, err = Init(endpoints,
		WithDialTimeout(time.Second),
		WithSecure("foo", "notfound.crt"))
	assert.Error(t, err)
	endpoints = nil
	_, err = Init(endpoints)
	assert.Error(t, err)
}
