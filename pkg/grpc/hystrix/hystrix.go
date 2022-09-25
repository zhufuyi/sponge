package hystrix

import (
	"context"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	metricCollector "github.com/afex/hystrix-go/hystrix/metric_collector"
	"github.com/afex/hystrix-go/plugins"
	"google.golang.org/grpc"
)

// https://github.com/soyacen/grpc-middleware/tree/main/hystrix

// UnaryClientInterceptor set the hystrix of unary client interceptor
func UnaryClientInterceptor(commandName string, opts ...Option) grpc.UnaryClientInterceptor {
	o := defaultOptions()
	o.apply(opts...)

	if o.statsD != nil {
		c, err := plugins.InitializeStatsdCollector(o.statsD)
		if err != nil {
			panic(err)
		}
		metricCollector.Registry.Register(c.NewStatsdCollector)
	}

	hystrix.ConfigureCommand(commandName, hystrix.CommandConfig{
		Timeout:                durationToInt(o.timeout, time.Millisecond),
		MaxConcurrentRequests:  o.maxConcurrentRequests,
		RequestVolumeThreshold: o.requestVolumeThreshold,
		SleepWindow:            durationToInt(o.sleepWindow, time.Millisecond),
		ErrorPercentThreshold:  o.errorPercentThreshold,
	})

	return unaryClientInterceptor(commandName, o)
}

func unaryClientInterceptor(commandName string, o *options) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		err = hystrix.DoC(ctx, commandName,
			// 熔断开关
			func(ctx context.Context) error {
				err = invoker(ctx, method, req, reply, cc, opts...)
				if err != nil {
					return err
				}
				return nil
			},
			// 降级处理
			o.fallbackFunc,
		)
		return err
	}
}

// StreamClientInterceptor set the hystrix of stream client interceptor
func StreamClientInterceptor(commandName string, opts ...Option) grpc.StreamClientInterceptor {
	o := defaultOptions()
	o.apply(opts...)

	if o.statsD != nil {
		c, err := plugins.InitializeStatsdCollector(o.statsD)
		if err != nil {
			panic(err)
		}
		metricCollector.Registry.Register(c.NewStatsdCollector)
	}

	hystrix.ConfigureCommand(commandName, hystrix.CommandConfig{
		Timeout:                durationToInt(o.timeout, time.Millisecond),
		MaxConcurrentRequests:  o.maxConcurrentRequests,
		RequestVolumeThreshold: o.requestVolumeThreshold,
		SleepWindow:            durationToInt(o.sleepWindow, time.Millisecond),
		ErrorPercentThreshold:  o.errorPercentThreshold,
	})

	return streamClientInterceptor(commandName, o)
}

func streamClientInterceptor(commandName string, o *options) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		var clientStream grpc.ClientStream
		err := hystrix.DoC(ctx, commandName,
			// 熔断开关
			func(ctx context.Context) error {
				var err error
				clientStream, err = streamer(ctx, desc, cc, method, opts...)
				if err != nil {
					return err
				}
				return nil
			},
			// 降级处理
			o.fallbackFunc,
		)
		return clientStream, err
	}
}

func durationToInt(duration time.Duration, unit time.Duration) int {
	durationAsNumber := duration / unit
	if int64(durationAsNumber) > int64(maxInt) {
		return maxInt
	}
	return int(durationAsNumber)
}
