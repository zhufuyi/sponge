// Package ws provides a websocket server implementation.
package ws

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// ServerOption is a functional option for the Server.
type ServerOption func(*serverOptions)

type serverOptions struct {
	responseHeader      http.Header
	upgrader            *websocket.Upgrader
	noClientPingTimeout time.Duration
	zapLogger           *zap.Logger
}

func defaultServerOptions() *serverOptions {
	return &serverOptions{
		upgrader: &websocket.Upgrader{ // default upgrader
			CheckOrigin: func(r *http.Request) bool { // allow all origins
				return true
			},
		},
	}
}

func (o *serverOptions) apply(opts ...ServerOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithResponseHeader sets the response header for the WebSocket upgrade response.
func WithResponseHeader(header http.Header) ServerOption {
	return func(o *serverOptions) {
		o.responseHeader = header
	}
}

// WithUpgrader sets the WebSocket upgrader for the server.
func WithUpgrader(upgrader *websocket.Upgrader) ServerOption {
	return func(o *serverOptions) {
		o.upgrader = upgrader
	}
}

// WithMaxMessageWaitPeriod sets the maximum waiting period for a message before closing the connection.
// Deprecated: use WithNoClientPingTimeout instead.
func WithMaxMessageWaitPeriod(period time.Duration) ServerOption {
	return func(o *serverOptions) {
		o.noClientPingTimeout = period
	}
}

// WithNoClientPingTimeout sets the timeout for the client to send a ping message, if timeout, the connection will be closed.
func WithNoClientPingTimeout(timeout time.Duration) ServerOption {
	return func(o *serverOptions) {
		o.noClientPingTimeout = timeout
	}
}

// WithServerLogger sets the logger for the server.
func WithServerLogger(l *zap.Logger) ServerOption {
	return func(o *serverOptions) {
		if l != nil {
			o.zapLogger = l
		}
	}
}

// --------------------------------------------------------------------------------------

// Conn is a WebSocket connection.
type Conn = websocket.Conn

// LoopFn is the function that is called for each WebSocket connection.
type LoopFn func(ctx context.Context, conn *Conn)

// Server is a WebSocket server.
type Server struct {
	upgrader *websocket.Upgrader

	w              http.ResponseWriter
	r              *http.Request
	responseHeader http.Header

	// If it is greater than 0, it means that the message waiting timeout mechanism is enabled
	//and the connection will be closed after the timeout, if it is 0, it means that the message
	// waiting timeout mechanism is not enabled.
	noClientPingTimeout time.Duration

	loopFn LoopFn

	zapLogger *zap.Logger
}

// NewServer creates a new WebSocket server.
func NewServer(w http.ResponseWriter, r *http.Request, loopFn LoopFn, opts ...ServerOption) *Server {
	o := defaultServerOptions()
	o.apply(opts...)
	if o.zapLogger == nil {
		o.zapLogger, _ = zap.NewProduction()
	}

	return &Server{
		w:      w,
		r:      r,
		loopFn: loopFn,

		upgrader:            o.upgrader,
		responseHeader:      o.responseHeader,
		noClientPingTimeout: o.noClientPingTimeout,
		zapLogger:           o.zapLogger,
	}
}

// Run runs the WebSocket server.
func (s *Server) Run(ctx context.Context) error {
	conn, err := s.upgrader.Upgrade(s.w, s.r, s.responseHeader)
	if err != nil {
		return err
	}
	defer conn.Close() //nolint

	fields := []zap.Field{zap.String("client", conn.RemoteAddr().String())}
	if s.noClientPingTimeout > 0 {
		// Set initial read deadline
		if err = conn.SetReadDeadline(time.Now().Add(s.noClientPingTimeout)); err != nil {
			return err
		}

		// Set up Ping handling for the connection,
		// when the client sends a ping message, the server side triggers this callback function
		conn.SetPingHandler(func(string) error {
			return conn.SetReadDeadline(time.Now().Add(s.noClientPingTimeout))
		})
		fields = append(fields, zap.String("no_ping_timeout", fmt.Sprintf("%vs", s.noClientPingTimeout.Seconds())))
	}

	s.zapLogger.Info("new websocket connection established", fields...)

	s.loopFn(ctx, conn)

	return nil
}

// IsClientClose returns true if the error is caused by client close.
func IsClientClose(err error) bool {
	return strings.Contains(err.Error(), "websocket: close") ||
		strings.Contains(err.Error(), "closed by the remote host")
}
