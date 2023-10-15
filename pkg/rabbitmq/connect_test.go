package rabbitmq

import (
	"context"
	"crypto/tls"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/utils"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	url    = "amqp://guest:guest@192.168.3.37:5672/"
	urlTLS = "amqps://guest:guest@127.0.0.1:5672/"
)

func TestConnectionOptions(t *testing.T) {
	opts := []ConnectionOption{
		WithLogger(nil),
		WithLogger(zap.NewNop()),
		WithReconnectTime(time.Second),
		WithTLSConfig(nil),
		WithTLSConfig(&tls.Config{
			InsecureSkipVerify: true,
		}),
	}

	o := defaultConnectionOptions()
	o.apply(opts...)

}

func TestNewConnection1(t *testing.T) {
	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		c, err := NewConnection("")
		assert.Error(t, err)

		c, err = NewConnection(url)
		if err != nil {
			t.Log(err)
			return
		}
		assert.True(t, c.CheckConnected())
		time.Sleep(time.Second)
		c.Close()

	})
}

func TestNewConnection2(t *testing.T) {
	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {

		// error
		_, err := NewConnection(urlTLS)
		assert.Error(t, err)

		_, err = NewConnection(urlTLS, WithTLSConfig(&tls.Config{
			InsecureSkipVerify: true,
		}))
		assert.Error(t, err)

	})
}

func TestConnection_monitor(t *testing.T) {
	c := &Connection{
		url:           urlTLS,
		reconnectTime: time.Second,
		Exit:          make(chan struct{}),
		ZapLog:        defaultLogger,
		Conn:          &amqp.Connection{},
		blockChan:     make(chan amqp.Blocking, 1),
		closeChan:     make(chan *amqp.Error, 1),
		IsConnected:   true,
	}

	c.CheckConnected()
	go func() {
		defer func() { recover() }()
		c.monitor()
	}()

	time.Sleep(time.Millisecond * 500)
	c.Mutex.Lock()
	c.blockChan <- amqp.Blocking{Active: false}
	c.blockChan <- amqp.Blocking{Active: true, Reason: "the disk is full."}
	c.Mutex.Unlock()

	time.Sleep(time.Millisecond * 500)
	c.Mutex.Lock()
	c.closeChan <- &amqp.Error{Code: 504, Reason: "connect failed"}
	c.Mutex.Unlock()

	time.Sleep(time.Millisecond * 500)
	c.Close()
	time.Sleep(time.Millisecond * 500)
}
