package initial

import (
	"context"
	"time"

	"github.com/zhufuyi/sponge/internal/config"
	"github.com/zhufuyi/sponge/internal/rpcclient"

	"github.com/zhufuyi/sponge/pkg/app"
	"github.com/zhufuyi/sponge/pkg/tracer"
)

// RegisterClose register for released resources
func RegisterClose(servers []app.IServer) []app.Close {
	var closes []app.Close

	// close server
	for _, s := range servers {
		closes = append(closes, s.Stop)
	}

	// close the rpc client connection
	closes = append(closes, func() error {
		return rpcclient.CloseServerNameExampleRPCConn()
	})

	// close tracing
	if config.Get().App.EnableTracing {
		closes = append(closes, func() error {
			ctx, _ := context.WithTimeout(context.Background(), 2*time.Second) //nolint
			return tracer.Close(ctx)
		})
	}

	return closes
}
