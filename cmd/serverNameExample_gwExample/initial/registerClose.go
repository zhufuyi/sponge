package initial

import (
	"context"
	"time"

	"github.com/zhufuyi/sponge/internal/config"
	"github.com/zhufuyi/sponge/internal/rpcclient"

	"github.com/zhufuyi/sponge/pkg/app"
	"github.com/zhufuyi/sponge/pkg/tracer"
)

// RegisterClose 注册app需要释放的资源
func RegisterClose(servers []app.IServer) []app.Close {
	var closes []app.Close

	// 关闭服务
	for _, s := range servers {
		closes = append(closes, s.Stop)
	}

	// 关闭rpc client连接
	closes = append(closes, func() error {
		return rpcclient.CloseServerNameExampleRPCConn()
	})

	// 关闭trace
	if config.Get().App.EnableTracing {
		closes = append(closes, func() error {
			ctx, _ := context.WithTimeout(context.Background(), 2*time.Second) //nolint
			return tracer.Close(ctx)
		})
	}

	return closes
}
