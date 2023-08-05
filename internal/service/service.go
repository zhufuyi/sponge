// Package service A grpc server-side or client-side package that handles business logic.
package service

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthPB "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	// registerFns collection of registration methods
	registerFns []func(server *grpc.Server)
)

// RegisterAllService register all services to the service
func RegisterAllService(server *grpc.Server) {
	healthPB.RegisterHealthServer(server, health.NewServer()) // Register for Health Screening

	for _, fn := range registerFns {
		fn(server)
	}
}
