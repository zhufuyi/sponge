package service

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthPB "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	// registerFns 注册方法集合
	registerFns []func(server *grpc.Server)
)

// RegisterAllService 注册所有service到服务中
func RegisterAllService(server *grpc.Server) {
	healthPB.RegisterHealthServer(server, health.NewServer()) // 注册健康检测

	for _, fn := range registerFns {
		fn(server)
	}
}
