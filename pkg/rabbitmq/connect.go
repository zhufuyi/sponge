// Package rabbitmq is a go wrapper for rabbitmq
package rabbitmq

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// Connection rabbitmq connection
type Connection struct {
	Mutex sync.Mutex

	url           string
	tlsConfig     *tls.Config
	reconnectTime time.Duration
	Exit          chan struct{}
	ZapLog        *zap.Logger

	Conn        *amqp.Connection
	blockChan   chan amqp.Blocking
	closeChan   chan *amqp.Error
	IsConnected bool
}

// NewConnection rabbitmq connection
func NewConnection(url string, opts ...ConnectionOption) (*Connection, error) {
	if url == "" {
		return nil, errors.New("url is empty")
	}

	o := defaultConnectionOptions()
	o.apply(opts...)

	c := &Connection{
		url:           url,
		reconnectTime: o.reconnectTime,
		tlsConfig:     o.tlsConfig,
		Exit:          make(chan struct{}),
		ZapLog:        o.zapLog,
	}

	conn, err := connect(c.url, c.tlsConfig)
	if err != nil {
		return nil, err
	}

	c.Conn = conn
	c.blockChan = c.Conn.NotifyBlocked(make(chan amqp.Blocking, 1))
	c.closeChan = c.Conn.NotifyClose(make(chan *amqp.Error, 1))
	c.IsConnected = true

	go c.monitor()

	return c, nil
}

func connect(url string, tlsConfig *tls.Config) (*amqp.Connection, error) {
	var (
		conn *amqp.Connection
		err  error
	)

	if strings.HasPrefix(url, "amqps://") {
		if tlsConfig == nil {
			return nil, errors.New("tls not set, e.g. NewConnection(url, WithTLSConfig(tlsConfig))")
		}
		conn, err = amqp.DialTLS(url, tlsConfig)
		if err != nil {
			return nil, err
		}
	} else {
		conn, err = amqp.Dial(url)
		if err != nil {
			return nil, err
		}
	}

	return conn, nil
}

// CheckConnected rabbitmq connection
func (c *Connection) CheckConnected() bool {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	return c.IsConnected
}

func (c *Connection) monitor() {
	retryCount := 0
	reconnectTip := fmt.Sprintf("[rabbitmq connection] lost connection, attempting reconnect in %s", c.reconnectTime)

	for {
		select {
		case <-c.Exit:
			_ = c.closeConn()
			c.ZapLog.Info("[rabbitmq connection] close connection")
			return
		case b := <-c.blockChan:
			if b.Active {
				c.ZapLog.Warn("[rabbitmq connection] TCP blocked: " + b.Reason)
			} else {
				c.ZapLog.Warn("[rabbitmq connection] TCP unblocked")
			}
		case <-c.closeChan:
			c.Mutex.Lock()
			c.IsConnected = false
			c.Mutex.Unlock()

			retryCount++
			c.ZapLog.Warn(reconnectTip)
			time.Sleep(c.reconnectTime) // wait for reconnect

			amqpConn, amqpErr := connect(c.url, c.tlsConfig)
			if amqpErr != nil {
				c.ZapLog.Warn("[rabbitmq connection] reconnect failed", zap.String("err", amqpErr.Error()), zap.Int("retryCount", retryCount))
				continue
			}
			c.ZapLog.Info("[rabbitmq connection] reconnect success")

			// set new connection
			c.Mutex.Lock()
			c.IsConnected = true
			c.Conn = amqpConn
			c.blockChan = c.Conn.NotifyBlocked(make(chan amqp.Blocking, 1))
			c.closeChan = c.Conn.NotifyClose(make(chan *amqp.Error, 1))
			c.Mutex.Unlock()
		}
	}
}

// Close rabbitmq connection
func (c *Connection) Close() {
	c.Mutex.Lock()
	c.IsConnected = false
	c.Mutex.Unlock()

	close(c.Exit)
}

func (c *Connection) closeConn() error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if c.Conn != nil {
		return c.Conn.Close()
	}

	return nil
}
