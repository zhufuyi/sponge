// Package gotest is a library that simulates the testing of cache, dao and handler.
package gotest

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zhufuyi/sponge/pkg/utils"
)

// Service info
type Service struct {
	Ctx      context.Context
	TestData interface{}
	MockDao  *Dao

	Server *grpc.Server
	listen net.Listener

	clientAddr     string
	clientConn     *grpc.ClientConn
	IServiceClient interface{}
}

// NewService instantiated service
func NewService(dao *Dao, testData interface{}) *Service {
	port, _ := utils.GetAvailablePort()
	clientAddr := fmt.Sprintf("127.0.0.1:%d", port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()

	return &Service{
		Ctx:      context.Background(),
		TestData: testData,
		MockDao:  dao,

		clientAddr: clientAddr,
		Server:     server,
		listen:     lis,
	}
}

// GoGrpcServer run grpc server
func (s *Service) GoGrpcServer() {
	go func() {
		if err := s.Server.Serve(s.listen); err != nil {
			panic(err)
		}
	}()
}

// GetClientConn dial rpc server
func (s *Service) GetClientConn() *grpc.ClientConn {
	conn, err := grpc.Dial(s.clientAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	return conn
}

// Close service
func (s *Service) Close() {
	if s.MockDao != nil {
		s.MockDao.Close()
	}
	if s.clientConn != nil {
		_ = s.clientConn.Close()
	}
	if s.Server != nil {
		s.Server.GracefulStop()
	}
}
