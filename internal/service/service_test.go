package service

import (
	"context"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/utils"

	"google.golang.org/grpc"
)

func TestRegisterAllService(t *testing.T) {
	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		server := grpc.NewServer()
		RegisterAllService(server)
		cancel()
	})
}
