package initial

import (
	"context"
	"time"

	"github.com/zhufuyi/sponge/internal/config"
	//"github.com/zhufuyi/sponge/internal/model"

	"github.com/zhufuyi/sponge/pkg/app"
	"github.com/zhufuyi/sponge/pkg/tracer"
)

// Close releasing resources after service exit
func Close(servers []app.IServer) []app.Close {
	var closes []app.Close

	// close server
	for _, s := range servers {
		closes = append(closes, s.Stop)
	}

	// close database
	//closes = append(closes, func() error {
	//	return model.CloseDB()
	//})

	// close redis
	//if config.Get().App.CacheType == "redis" {
	//	closes = append(closes, func() error {
	//		return model.CloseRedis()
	//	})
	//}

	// close tracing
	if config.Get().App.EnableTrace {
		closes = append(closes, func() error {
			ctx, _ := context.WithTimeout(context.Background(), 2*time.Second) //nolint
			return tracer.Close(ctx)
		})
	}

	return closes
}
