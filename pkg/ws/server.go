// Package ws provides a WebSocket server implementation.
package ws

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// ServerOption is a functional option for the Server.
type ServerOption func(*serverOptions)

type serverOptions struct {
	responseHeader       http.Header
	upgrader             *websocket.Upgrader
	maxMessageWaitPeriod time.Duration
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
func WithMaxMessageWaitPeriod(period time.Duration) ServerOption {
	return func(o *serverOptions) {
		o.maxMessageWaitPeriod = period
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
	maxMessageWaitPeriod time.Duration

	loopFn LoopFn
}

// NewServer creates a new WebSocket server.
func NewServer(w http.ResponseWriter, r *http.Request, loopFn LoopFn, opts ...ServerOption) *Server {
	o := defaultServerOptions()
	o.apply(opts...)

	return &Server{
		w:      w,
		r:      r,
		loopFn: loopFn,

		upgrader:             o.upgrader,
		responseHeader:       o.responseHeader,
		maxMessageWaitPeriod: o.maxMessageWaitPeriod,
	}
}

// Run runs the WebSocket server.
func (s *Server) Run(ctx context.Context) error {
	conn, err := s.upgrader.Upgrade(s.w, s.r, s.responseHeader)
	if err != nil {
		return err
	}
	defer conn.Close() //nolint

	if s.maxMessageWaitPeriod > 0 {
		// Set initial read deadline
		err = conn.SetReadDeadline(time.Now().Add(s.maxMessageWaitPeriod))
		if err != nil {
			return err
		}

		// Set up Ping handling for the connection,
		// when the client sends a ping message, the server side triggers this callback function
		conn.SetPingHandler(func(string) error {
			_ = conn.SetReadDeadline(time.Now().Add(s.maxMessageWaitPeriod))
			return nil
		})
	}

	s.loopFn(ctx, conn)

	return nil
}

// IsClientClose returns true if the error is caused by client close.
func IsClientClose(err error) bool {
	return strings.Contains(err.Error(), "websocket: close") ||
		strings.Contains(err.Error(), "closed by the remote host")
}
