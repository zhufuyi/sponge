package interceptor

import (
	"time"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

// ---------------------------------- client interceptor ----------------------------------

// UnaryClientLog 客户端日志unary拦截器
func UnaryClientLog(logger *zap.Logger, opts ...grpc_zap.Option) grpc.UnaryClientInterceptor {
	return grpc_zap.UnaryClientInterceptor(logger, opts...)
}

// UnaryStreamLog 客户端日志stream拦截器
func UnaryStreamLog(logger *zap.Logger, opts ...grpc_zap.Option) grpc.StreamClientInterceptor {
	return grpc_zap.StreamClientInterceptor(logger, opts...)
}

// ---------------------------------- server interceptor ----------------------------------

var ignoreLogMethods = map[string]struct{}{} // 忽略打印的方法

// LogOption 日志设置
type LogOption func(*logOptions)

type logOptions struct {
	fields        map[string]interface{}
	ignoreMethods map[string]struct{}
}

func defaultLogOptions() *logOptions {
	return &logOptions{
		fields:        make(map[string]interface{}), // 自定义打印kv
		ignoreMethods: make(map[string]struct{}),    // 忽略打印日志的方法
	}
}

func (o *logOptions) apply(opts ...LogOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithLogFields 添加自定义打印字段
func WithLogFields(kvs map[string]interface{}) LogOption {
	return func(o *logOptions) {
		if len(kvs) == 0 {
			return
		}
		o.fields = kvs
	}
}

// WithLogIgnoreMethods 忽略打印的方法
// fullMethodName格式: /packageName.serviceName/methodName，
// 示例/api.userExample.v1.userExampleService/GetByID
func WithLogIgnoreMethods(fullMethodNames ...string) LogOption {
	return func(o *logOptions) {
		for _, method := range fullMethodNames {
			o.ignoreMethods[method] = struct{}{}
		}
	}
}

// UnaryServerLog 服务端日志unary拦截器
func UnaryServerLog(logger *zap.Logger, opts ...LogOption) grpc.UnaryServerInterceptor {
	o := defaultLogOptions()
	o.apply(opts...)
	ignoreLogMethods = o.ignoreMethods

	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	grpc_zap.ReplaceGrpcLoggerV2(logger)

	// 日志设置，默认打印客户端断开连接信息，示例 https://pkg.go.dev/github.com/grpc-ecosystem/go-grpc-middleware/logging/zap
	zapOptions := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_us", duration.Microseconds()) // 默认打印耗时字段
		}),
	}

	// 自定义打印字段
	for key, val := range o.fields {
		zapOptions = append(zapOptions, grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Any(key, val)
		}))
	}

	// 自定义跳过打印日志的调用方法
	if len(ignoreLogMethods) > 0 {
		zapOptions = append(zapOptions, grpc_zap.WithDecider(func(fullMethodName string, err error) bool {
			if err == nil {
				if _, ok := ignoreLogMethods[fullMethodName]; ok {
					return false
				}
			}
			return true
		}))
	}

	return grpc_zap.UnaryServerInterceptor(logger, zapOptions...)
}

// UnaryServerCtxTags extractor field unary拦截器
func UnaryServerCtxTags() grpc.UnaryServerInterceptor {
	return grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor))
}

// StreamServerLog 服务端日志stream拦截器
func StreamServerLog(logger *zap.Logger, opts ...LogOption) grpc.StreamServerInterceptor {
	o := defaultLogOptions()
	o.apply(opts...)
	ignoreLogMethods = o.ignoreMethods

	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	grpc_zap.ReplaceGrpcLoggerV2(logger)

	// 日志设置，默认打印客户端断开连接信息，示例 https://pkg.go.dev/github.com/grpc-ecosystem/go-grpc-middleware/logging/zap
	zapOptions := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_us", duration.Microseconds()) // 默认打印耗时字段
		}),
	}

	// 自定义打印字段
	for key, val := range o.fields {
		zapOptions = append(zapOptions, grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Any(key, val)
		}))
	}

	// 自定义跳过打印日志的调用方法
	if len(ignoreLogMethods) > 0 {
		zapOptions = append(zapOptions, grpc_zap.WithDecider(func(fullMethodName string, err error) bool {
			if err == nil {
				if _, ok := ignoreLogMethods[fullMethodName]; ok {
					return false
				}
			}
			return true
		}))
	}

	return grpc_zap.StreamServerInterceptor(logger, zapOptions...)
}

// StreamServerCtxTags extractor field stream拦截器
func StreamServerCtxTags() grpc.StreamServerInterceptor {
	return grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor))
}
