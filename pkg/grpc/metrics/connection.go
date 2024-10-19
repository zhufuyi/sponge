package metrics

import (
	"fmt"
	"net"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// ConnectionOption set connection option
type ConnectionOption func(*connectionOptions)

type connectionOptions struct {
	zapLogger       *zap.Logger
	connectionGauge prometheus.Gauge
}

func defaultConnectionOptions() *connectionOptions {
	return &connectionOptions{}
}

func (o *connectionOptions) apply(opts ...ConnectionOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithConnectionsLogger set logger for connection
func WithConnectionsLogger(l *zap.Logger) ConnectionOption {
	return func(o *connectionOptions) {
		if l != nil {
			o.zapLogger = l
		}
	}
}

// WithConnectionsGauge set prometheus gauge for connections
func WithConnectionsGauge() ConnectionOption {
	return func(o *connectionOptions) {
		o.connectionGauge = grpcConnectionGauge
	}
}

// ------------------------------------------------------------------------------------------

// CustomConn custom connections, intercept disconnected behavior
type CustomConn struct {
	net.Conn
	listener *CustomListener
}

// CustomListener custom listener for counting connections
type CustomListener struct {
	net.Listener
	activeConnections int
	mu                sync.Mutex
	zapLogger         *zap.Logger
	connectionGauge   prometheus.Gauge
}

// Accept waits for and returns the next connection to the listener.
func (l *CustomListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	var count int
	l.mu.Lock()
	l.activeConnections++
	count = l.activeConnections
	l.mu.Unlock()

	if l.zapLogger != nil {
		l.zapLogger.Info("new grpc client connected", zap.String("client", conn.RemoteAddr().String()), zap.Int("active connections", count))
	} else {
		fmt.Printf("new grpc client connected, client: %s, active connections: %d\n", conn.RemoteAddr().String(), count)
	}
	if l.connectionGauge != nil {
		l.connectionGauge.Set(float64(count))
	}

	return &CustomConn{
		Conn:     conn,
		listener: l,
	}, nil
}

// GetActiveConnections returns the number of active connections.
func (l *CustomListener) GetActiveConnections() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.activeConnections
}

// closes the connection and decrements the active connections count.
func (l *CustomListener) closeConnection(clientAddr string) {
	var count int
	l.mu.Lock()
	l.activeConnections--
	count = l.activeConnections
	l.mu.Unlock()

	if l.zapLogger != nil {
		l.zapLogger.Info("grpc client disconnected", zap.String("client", clientAddr), zap.Int("active connections", count))
	} else {
		fmt.Printf("grpc client disconnected client: %s, active connections: %d\n", clientAddr, count)
	}
	if l.connectionGauge != nil {
		l.connectionGauge.Set(float64(count))
	}
}

// Close closes the listener, any blocked except operations will be unblocked and return errors.
func (c *CustomConn) Close() error {
	defer func() { _ = recover() }()
	clientAddr := c.Conn.RemoteAddr().String()
	err := c.Conn.Close()
	if err == nil {
		c.listener.closeConnection(clientAddr)
	}
	return err
}

// NewCustomListener creates a new custom listener.
func NewCustomListener(listener net.Listener, opts ...ConnectionOption) *CustomListener {
	o := defaultConnectionOptions()
	o.apply(opts...)

	return &CustomListener{
		Listener:        listener,
		zapLogger:       o.zapLogger,
		connectionGauge: o.connectionGauge,
	}
}
