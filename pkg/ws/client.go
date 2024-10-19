package ws

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
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

	zapLogger *zap.Logger
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

// WithClientLogger sets the logger for the client.
func WithClientLogger(l *zap.Logger) ClientOption {
	return func(o *clientOptions) {
		if l != nil {
			o.zapLogger = l
		}
	}
}

// ----------------------------------------------------------------------------------

// Client is a wrapper of gorilla/websocket.
type Client struct {
	dialer        *websocket.Dialer
	requestHeader http.Header
	url           string
	conn          *websocket.Conn

	pingInterval time.Duration
	ctx          context.Context
	cancel       context.CancelFunc
	zapLogger    *zap.Logger
}

// NewClient creates a new client.
func NewClient(url string, opts ...ClientOption) (*Client, error) {
	o := defaultClientOptions()
	o.apply(opts...)
	if o.zapLogger == nil {
		o.zapLogger, _ = zap.NewProduction()
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := &Client{
		url:           url,
		dialer:        o.dialer,
		requestHeader: o.requestHeader,
		pingInterval:  o.pingDialInterval,
		ctx:           ctx,
		cancel:        cancel,
		zapLogger:     o.zapLogger,
	}

	err := c.connect()
	if err != nil {
		return nil, err
	}

	fields := []zap.Field{zap.String("server", c.url)}
	if c.pingInterval > 0 {
		c.ping()
		fields = append(fields, zap.String("auto ping interval", fmt.Sprintf("%vs", c.pingInterval.Seconds())))
	}

	c.zapLogger.Info("connect websocket server success", fields...)

	return c, nil
}

// GetConn returns the connection of the client.
func (c *Client) GetConn() *websocket.Conn {
	if c.conn == nil {
		defer func() {
			if e := recover(); e != nil {
				c.zapLogger.Warn("connect websocket server error", zap.Any("err", e))
			}
		}()
		err := c.connect()
		if err != nil {
			panic(err)
		}
	}

	return c.conn
}

// connect the websocket server.
func (c *Client) connect() error {
	conn, _, err := c.dialer.Dial(c.url, c.requestHeader)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

// TryReconnect tries to reconnect the websocket server.
func (c *Client) TryReconnect() error {
	delay := 1 * time.Second
	maxDelay := 32 * time.Second
	for {
		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		case <-time.After(delay):
			if err := c.connect(); err != nil {
				if delay >= maxDelay {
					delay = maxDelay
					c.zapLogger.Warn("reconnect websocket server error", zap.Error(err), zap.String("server", c.url))
					continue
				}
				delay *= 2
				continue
			}
			c.zapLogger.Warn("reconnect websocket server success", zap.String("server", c.url))
			return nil
		}
	}
}

// ping websocket server, try to reconnect if connection failed.
func (c *Client) ping() {
	go func() {
		isExit := false
		defer func() {
			if e := recover(); e != nil {
				c.zapLogger.Warn("ping server panic", zap.Any("err", e))
			}

			if !isExit {
				if err := c.TryReconnect(); err == nil {
					c.ping()
				}
			}
		}()

		ticker := time.NewTicker(c.pingInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := c.conn.WriteControl(websocket.PingMessage, pingData, time.Now().Add(5*time.Second)); err != nil {
					c.zapLogger.Warn("ping server error", zap.Error(err))
					return
				}

			case <-c.ctx.Done(): // exit
				isExit = true
				return
			}
		}
	}()
}

// GetCtx returns the context of the client.
func (c *Client) GetCtx() context.Context {
	return c.ctx
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
