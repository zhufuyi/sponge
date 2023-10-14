package rabbitmq

import (
	"crypto/tls"
	"time"

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
