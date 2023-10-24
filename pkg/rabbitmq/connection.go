// Package rabbitmq is a go wrapper for github.com/rabbitmq/amqp091-go
//
// producer and consumer using the five types direct, topic, fanout, headers, x-delayed-message.
// publisher and subscriber using the fanout message type.
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

// DefaultURL default rabbitmq url
const DefaultURL = "amqp://guest:guest@localhost:5672/"

var defaultLogger, _ = zap.NewProduction()

// ConnectionOption connection option.
type ConnectionOption func(*connectionOptions)

type connectionOptions struct {
	tlsConfig     *tls.Config   // tls config, if the url is amqps this field must be set
	reconnectTime time.Duration // reconnect time interval, default is 3s

	zapLog *zap.Logger
}

func (o *connectionOptions) apply(opts ...ConnectionOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// default connection settings
func defaultConnectionOptions() *connectionOptions {
	return &connectionOptions{
		tlsConfig:     nil,
		reconnectTime: time.Second * 3,
		zapLog:        defaultLogger,
	}
}

// WithTLSConfig set tls config option.
func WithTLSConfig(tlsConfig *tls.Config) ConnectionOption {
	return func(o *connectionOptions) {
		if tlsConfig == nil {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		o.tlsConfig = tlsConfig
	}
}

// WithReconnectTime set reconnect time interval option.
func WithReconnectTime(d time.Duration) ConnectionOption {
	return func(o *connectionOptions) {
		if d == 0 {
			d = time.Second * 3
		}
		o.reconnectTime = d
	}
}

// WithLogger set logger option.
func WithLogger(zapLog *zap.Logger) ConnectionOption {
	return func(o *connectionOptions) {
		if zapLog == nil {
			return
		}
		o.zapLog = zapLog
	}
}

// -------------------------------------------------------------------------------------------

// Connection rabbitmq connection
type Connection struct {
	mutex sync.Mutex

	url           string
	tlsConfig     *tls.Config
	reconnectTime time.Duration
	exit          chan struct{}
	zapLog        *zap.Logger

	conn        *amqp.Connection
	blockChan   chan amqp.Blocking
	closeChan   chan *amqp.Error
	isConnected bool
}

// NewConnection rabbitmq connection
func NewConnection(url string, opts ...ConnectionOption) (*Connection, error) {
	if url == "" {
		return nil, errors.New("url is empty")
	}

	o := defaultConnectionOptions()
	o.apply(opts...)

	connection := &Connection{
		url:           url,
		reconnectTime: o.reconnectTime,
		tlsConfig:     o.tlsConfig,
		exit:          make(chan struct{}),
		zapLog:        o.zapLog,
	}

	conn, err := connect(connection.url, connection.tlsConfig)
	if err != nil {
		return nil, err
	}
	connection.zapLog.Info("[rabbitmq connection] connected successfully.")

	connection.conn = conn
	connection.blockChan = connection.conn.NotifyBlocked(make(chan amqp.Blocking, 1))
	connection.closeChan = connection.conn.NotifyClose(make(chan *amqp.Error, 1))
	connection.isConnected = true

	go connection.monitor()

	return connection, nil
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
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.isConnected
}

func (c *Connection) monitor() {
	retryCount := 0
	reconnectTip := fmt.Sprintf("[rabbitmq connection] lost connection, attempting reconnect in %s", c.reconnectTime)

	for {
		select {
		case <-c.exit:
			_ = c.closeConn()
			c.zapLog.Info("[rabbitmq connection] closed")
			return
		case b := <-c.blockChan:
			if b.Active {
				c.zapLog.Warn("[rabbitmq connection] TCP blocked: " + b.Reason)
			} else {
				c.zapLog.Warn("[rabbitmq connection] TCP unblocked")
			}
		case <-c.closeChan:
			c.mutex.Lock()
			c.isConnected = false
			c.mutex.Unlock()

			retryCount++
			c.zapLog.Warn(reconnectTip)
			time.Sleep(c.reconnectTime) // wait for reconnect

			amqpConn, amqpErr := connect(c.url, c.tlsConfig)
			if amqpErr != nil {
				c.zapLog.Warn("[rabbitmq connection] reconnect failed", zap.String("err", amqpErr.Error()), zap.Int("retryCount", retryCount))
				continue
			}
			c.zapLog.Info("[rabbitmq connection] reconnected successfully.")

			// set new connection
			c.mutex.Lock()
			c.isConnected = true
			c.conn = amqpConn
			c.blockChan = c.conn.NotifyBlocked(make(chan amqp.Blocking, 1))
			c.closeChan = c.conn.NotifyClose(make(chan *amqp.Error, 1))
			c.mutex.Unlock()
		}
	}
}

// Close rabbitmq connection
func (c *Connection) Close() {
	c.mutex.Lock()
	c.isConnected = false
	c.mutex.Unlock()

	close(c.exit)
}

func (c *Connection) closeConn() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}
