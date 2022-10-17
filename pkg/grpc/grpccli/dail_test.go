package grpccli

import (
	"context"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/registry/etcd"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc/credentials/insecure"
)

func TestDial(t *testing.T) {
	_, err := Dial(context.Background(), "localhost:8282")
	assert.NotNil(t, err)
}

func TestDialInsecure(t *testing.T) {
	_, err := DialInsecure(context.Background(), "localhost:8282")
	assert.NoError(t, err)
}

func Test_dial(t *testing.T) {
	_, err := dial(context.Background(), "localhost:8282", true,
		WithCredentials(insecure.NewCredentials()),
		WithEnableLog(zap.NewNop()),
		WithEnableMetrics(),
		WithEnableLoadBalance(),
		WithEnableCircuitBreaker(),
		WithEnableRetry(),
		WithDiscovery(etcd.New(&clientv3.Client{})),
	)
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 10)
}
