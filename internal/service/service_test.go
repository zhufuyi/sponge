package service

import (
	"testing"

	"google.golang.org/grpc"
)

func TestRegisterAllService(t *testing.T) {
	defer func() {
		recover()
	}()

	server := grpc.NewServer()
	RegisterAllService(server)
}
