// Package keepalive is setting grpc keepalive parameters.
package keepalive

import (
	"math"
	"time"

	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc"
)

// ---------------------------------- client option ----------------------------------

var kacp = keepalive.ClientParameters{
	Time:                20 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             1 * time.Second,  // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

// ClientKeepAlive keep the connection set
func ClientKeepAlive() grpc.DialOption {
	return grpc.WithKeepaliveParams(kacp)
}

// ---------------------------------- server option ----------------------------------

const (
	infinity                     = time.Duration(math.MaxInt64)
	defaultMaxConnectionIdle     = infinity
	defaultMaxConnectionAge      = infinity
	defaultMaxConnectionAgeGrace = infinity
)

var kaep = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
	PermitWithoutStream: true,            // Allow pings even when there are no active streams
}

var kasp = keepalive.ServerParameters{
	MaxConnectionIdle:     defaultMaxConnectionIdle,     // If a client is idle for 15 seconds, send a GOAWAY
	MaxConnectionAge:      defaultMaxConnectionAge,      // If any connection is alive for more than 30 seconds, send a GOAWAY
	MaxConnectionAgeGrace: defaultMaxConnectionAgeGrace, // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
	Time:                  20 * time.Second,             // Ping the client if it is idle for 5 seconds to ensure the connection is still active
	Timeout:               1 * time.Second,              // Wait 1 second for the ping ack before assuming the connection is dead
}

// ServerKeepAlive keep the connection set
func ServerKeepAlive() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
	}
}
