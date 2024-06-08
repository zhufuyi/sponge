package ws

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pingData = []byte("ping")
)

// ClientOption is a functional option for the client.
type ClientOption func(*clientOptions)

type clientOptions struct {
	dialer           *websocket.Dialer
	requestHeader    http.Header
	pingDialInterval time.Duration
}

func defaultClientOptions() *clientOptions {
	return &clientOptions{
		dialer: websocket.DefaultDialer,
	}
}

func (o *clientOptions) apply(opts ...ClientOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithDialer sets the dialer for the client.
func WithDialer(dialer *websocket.Dialer) ClientOption {
	return func(o *clientOptions) {
		o.dialer = dialer
	}
}

// WithRequestHeader sets the request header for the client.
func WithRequestHeader(header http.Header) ClientOption {
	return func(o *clientOptions) {
		o.requestHeader = header
	}
}

// WithPing sets the interval for sending ping message to the server.
func WithPing(interval time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.pingDialInterval = interval
	}
}

// ----------------------------------------------------------------------------------

// Client is a wrapper of gorilla/websocket.
type Client struct {
	dialer        *websocket.Dialer
	requestHeader http.Header
	url           string
	conn          *websocket.Conn

	pingDialInterval time.Duration
	ctx              context.Context
	cancel           context.CancelFunc

	once sync.Once
}

// NewClient creates a new client.
func NewClient(url string, opts ...ClientOption) (*Client, error) {
	o := defaultClientOptions()
	o.apply(opts...)

	ctx, cancel := context.WithCancel(context.Background())

	c := &Client{
		url:              url,
		dialer:           o.dialer,
		requestHeader:    o.requestHeader,
		pingDialInterval: o.pingDialInterval,
		ctx:              ctx,
		cancel:           cancel,
	}

	err := c.Reconnect()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) GetConn() *websocket.Conn {
	if c.conn == nil {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("get conn panic, %v\n", err)
			}
		}()
		err := c.Reconnect()
		if err != nil {
			panic(err)
		}
	}

	return c.conn
}

// Reconnect the websocket server.
func (c *Client) Reconnect() error {
	conn, _, err := c.dialer.Dial(c.url, c.requestHeader)
	if err != nil {
		return err
	}
	c.conn = conn

	if c.pingDialInterval > 0 {
		c.once.Do(func() {
			log.Println("start ping server")
			c.ping()
		})
	}

	return nil
}

// timed ping server
func (c *Client) ping() {
	go func() {
		isExit := false
		defer func() {
			if err := recover(); err != nil {
				log.Printf("ping server panic, %v\n", err)
			}

			if !isExit {
				time.Sleep(15 * time.Second)
				c.ping()
			}

			log.Printf("ping server exit\n")
		}()

		ticker := time.NewTicker(c.pingDialInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := c.conn.WriteControl(websocket.PingMessage, pingData, time.Now().Add(5*time.Second)); err != nil {
					log.Printf("ping server err, %v\n", err)
					continue
				}

			case <-c.ctx.Done():
				isExit = true
				return
			}
		}
	}()
}

// Close closes the connection.
// Note: if set pingDialInterval, the Close method must be called, otherwise it will cause the goroutine to leak
func (c *Client) Close() error {
	if c.cancel != nil {
		c.cancel()
	}
	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

// IsServerClose returns true if the error is caused by server close.
func IsServerClose(err error) bool {
	return strings.Contains(err.Error(), "use of closed network")
}
