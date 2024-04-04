package interceptor

import (
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/logger"
)

func TestUnaryClientLog(t *testing.T) {
	addr := newUnaryRPCServer()
	time.Sleep(time.Millisecond * 200)
	cli := newUnaryRPCClient(addr,
		UnaryClientRequestID(),
		UnaryClientLog(logger.Get(), WithReplaceGRPCLogger()),
	)
	_ = sayHelloMethod(cli)
}

func TestUnaryServerLog(t *testing.T) {
	addr := newUnaryRPCServer(
		UnaryServerRequestID(),
		UnaryServerLog(logger.Get(), WithReplaceGRPCLogger()),
		UnaryServerSimpleLog(logger.Get(), WithReplaceGRPCLogger()),
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
		StreamClientLog(logger.Get(), WithReplaceGRPCLogger()),
	)
	_ = discussHelloMethod(cli)
	time.Sleep(time.Millisecond)
}

func TestUnaryServerLog_ignore(t *testing.T) {
	addr := newUnaryRPCServer(
		UnaryServerLog(logger.Get(),
			WithLogFields(map[string]interface{}{"foo": "bar"}),
			WithLogIgnoreMethods("/api.user.v1.user/GetByID"),
		),
	)
	time.Sleep(time.Millisecond * 200)
	cli := newUnaryRPCClient(addr)
	_ = sayHelloMethod(cli)
}

func TestStreamServerLog(t *testing.T) {
	addr := newStreamRPCServer(
		StreamServerRequestID(),
		StreamServerLog(logger.Get(),
			WithReplaceGRPCLogger(),
			WithLogFields(map[string]interface{}{}),
		),
		StreamServerSimpleLog(logger.Get(),
			WithReplaceGRPCLogger(),
			WithLogFields(map[string]interface{}{}),
		),
	)
	time.Sleep(time.Millisecond * 200)
	cli := newStreamRPCClient(addr)
	_ = discussHelloMethod(cli)
	time.Sleep(time.Millisecond)
}

// ----------------------------------------------------------------------------------------

func TestNilLog(t *testing.T) {
	UnaryClientLog(nil)
	StreamClientLog(nil)
	UnaryServerLog(nil)
	UnaryServerSimpleLog(nil)
	StreamServerLog(nil)
	StreamServerSimpleLog(nil)
}
