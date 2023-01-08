package interceptor

import (
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/logger"

	"github.com/stretchr/testify/assert"
)

func TestUnaryClientLog(t *testing.T) {
	addr := newUnaryRPCServer()
	time.Sleep(time.Millisecond * 200)
	cli := newUnaryRPCClient(addr,
		UnaryClientRequestID(),
		UnaryClientLog(logger.Get()),
		UnaryClientLog2(logger.Get()),
	)
	_ = sayHelloMethod(cli)
}

func TestUnaryServerLog(t *testing.T) {
	addr := newUnaryRPCServer(
		UnaryServerRequestID(),
		UnaryServerLog(logger.Get()),
		UnaryServerLog2(
			logger.Get(),
			WithLogFields(map[string]interface{}{"foo": "bar"}),
			WithLogIgnoreMethods("/ping"),
		),
	)
	time.Sleep(time.Millisecond * 200)
	cli := newUnaryRPCClient(addr)
	_ = sayHelloMethod(cli)
}

func TestUnaryServerLog_ignore(t *testing.T) {
	addr := newUnaryRPCServer(
		UnaryServerLog(logger.Get()),
		UnaryServerLog2(
			logger.Get(),
			WithLogFields(nil),
			WithLogIgnoreMethods("/proto.Greeter/SayHello"),
		),
	)
	time.Sleep(time.Millisecond * 200)
	cli := newUnaryRPCClient(addr)
	_ = sayHelloMethod(cli)
}

func TestStreamClientLog(t *testing.T) {
	addr := newStreamRPCServer()
	time.Sleep(time.Millisecond * 200)
	cli := newStreamRPCClient(addr,
		StreamClientRequestID(),
		StreamClientLog(logger.Get()),
		StreamClientLog2(logger.Get()),
	)
	_ = discussHelloMethod(cli)
	time.Sleep(time.Millisecond)
}

func TestStreamServerLog(t *testing.T) {
	addr := newStreamRPCServer(
		StreamServerRequestID(),
		StreamServerLog(logger.Get()),
		StreamServerLog2(logger.Get()),
	)
	time.Sleep(time.Millisecond * 200)
	cli := newStreamRPCClient(addr)
	_ = discussHelloMethod(cli)
	time.Sleep(time.Millisecond)
}

// ----------------------------------------------------------------------------------------

func TestUnaryClientLog2(t *testing.T) {
	interceptor := UnaryClientLog2(nil)
	assert.NotNil(t, interceptor)
}

func TestStreamClientLog2(t *testing.T) {
	interceptor := StreamClientLog2(nil)
	assert.NotNil(t, interceptor)
}

func TestStreamServerLog2(t *testing.T) {
	interceptor := StreamServerLog2(nil,
		WithLogFields(map[string]interface{}{"foo": "bar"}),
		WithLogIgnoreMethods("/ping"),
	)
	assert.NotNil(t, interceptor)
}

func TestUnaryServerLog2(t *testing.T) {
	interceptor := UnaryServerLog2(nil,
		WithLogFields(map[string]interface{}{"foo": "bar"}),
		WithLogIgnoreMethods("/ping"),
	)
	assert.NotNil(t, interceptor)
}

func TestNilLog(t *testing.T) {
	UnaryClientLog(nil)
	UnaryClientLog2(nil)
	StreamClientLog(nil)
	StreamClientLog2(nil)
	UnaryServerLog(nil)
	UnaryServerLog2(nil)
	StreamServerLog(nil)
	StreamServerLog2(nil)
}
