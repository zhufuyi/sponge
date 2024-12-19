package grpccli

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"

	"github.com/go-dev-frame/sponge/pkg/grpc/gtls/certfile"
	"github.com/go-dev-frame/sponge/pkg/servicerd/registry/etcd"
)

func TestNewClient(t *testing.T) {
	_, err := NewClient("localhost:8282")
	assert.NoError(t, err)
	_, err = Dial(context.Background(), "localhost:8282")
	assert.NoError(t, err)
}

func TestNewClient2(t *testing.T) {
	_, err := NewClient("localhost:8282",
		WithEnableLog(zap.NewNop()),
		WithEnableMetrics(),
		WithToken(true, "grpc", "123456"),
		WithEnableLoadBalance(),
		WithEnableCircuitBreaker(),
		WithEnableRetry(),
		WithDiscovery(etcd.New(&clientv3.Client{})),
	)
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 50)
}

func Test_unaryClientOptions(t *testing.T) {
	o := &options{
		enableToken:          true,
		enableLog:            true,
		enableRequestID:      true,
		enableTrace:          true,
		enableMetrics:        true,
		enableRetry:          true,
		enableLoadBalance:    true,
		enableCircuitBreaker: true,
	}
	scOpt := unaryClientOptions(o)
	assert.NotNil(t, scOpt)
}

func Test_streamClientOptions(t *testing.T) {
	o := &options{
		enableToken:          true,
		enableLog:            true,
		enableRequestID:      true,
		enableTrace:          true,
		enableMetrics:        true,
		enableRetry:          true,
		enableLoadBalance:    true,
		enableCircuitBreaker: true,
	}
	scOpt := streamClientOptions(o)
	assert.NotNil(t, scOpt)
}

func Test_secureOption(t *testing.T) {
	o := &options{
		secureType: "one-way",
		serverName: "localhost",
		certFile:   certfile.Path("one-way/server.crt"),
	}

	// correct
	opt, err := secureOption(o)
	assert.NoError(t, err)
	assert.NotNil(t, opt)

	// error
	o.certFile = ""
	_, err = secureOption(o)
	assert.Error(t, err)
	o.certFile = "not found"
	_, err = secureOption(o)
	assert.Error(t, err)

	o = &options{
		secureType: "two-way",
		serverName: "localhost",
		caFile:     certfile.Path("two-way/ca.pem"),
		certFile:   certfile.Path("two-way/client/client.pem"),
		keyFile:    certfile.Path("two-way/client/client.key"),
	}

	// correct
	opt, err = secureOption(o)
	assert.NoError(t, err)
	assert.NotNil(t, opt)

	// error
	o.certFile = "not found"
	_, err = secureOption(o)
	assert.Error(t, err)
	o.caFile = ""
	_, err = secureOption(o)
	assert.Error(t, err)
	o.caFile = "not found"
	o.certFile = ""
	_, err = secureOption(o)
	assert.Error(t, err)
	o.certFile = "not found"
	o.keyFile = ""
	_, err = secureOption(o)
	assert.Error(t, err)

	o.secureType = ""
	opt, err = secureOption(o)
	assert.NoError(t, err)
	assert.NotNil(t, opt)
}
