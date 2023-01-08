package interceptor

import (
	"context"
	"encoding/json"
	"time"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// ---------------------------------- client interceptor ----------------------------------

// UnaryClientLog client log unary interceptor
func UnaryClientLog(logger *zap.Logger) grpc.UnaryClientInterceptor {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		startTime := time.Now()

		var reqIDField zap.Field
		if requestID := ClientCtxRequestID(ctx); requestID != "" {
			reqIDField = zap.String(ContextRequestIDKey, requestID)
		} else {
			reqIDField = zap.Skip()
		}

		err := invoker(ctx, method, req, reply, cc, opts...)

		fields := []zap.Field{
			zap.String("code", status.Code(err).String()),
			zap.Error(err),
			zap.String("rpc_type", "unary"),
			zap.String("method", method),
			zap.Int64("time_us", time.Since(startTime).Microseconds()),
			reqIDField,
		}

		logger.Info("rpc client invoker", fields...)
		return err
	}
}

// UnaryClientLog2 client log unary interceptor
func UnaryClientLog2(logger *zap.Logger, opts ...grpc_zap.Option) grpc.UnaryClientInterceptor {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	return grpc_zap.UnaryClientInterceptor(logger, opts...)
}

// StreamClientLog client log stream interceptor
func StreamClientLog(logger *zap.Logger) grpc.StreamClientInterceptor {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
		streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		startTime := time.Now()

		var reqIDField zap.Field
		if requestID := ClientCtxRequestID(ctx); requestID != "" {
			reqIDField = zap.String(ContextRequestIDKey, requestID)
		} else {
			reqIDField = zap.Skip()
		}

		clientStream, err := streamer(ctx, desc, cc, method, opts...)

		fields := []zap.Field{
			zap.String("code", status.Code(err).String()),
			zap.Error(err),
			zap.String("rpc_type", "stream"),
			zap.String("method", method),
			zap.Int64("time_us", time.Since(startTime).Microseconds()),
			reqIDField,
		}

		logger.Info("rpc client invoker", fields...)
		return clientStream, err
	}
}

// StreamClientLog2 client log stream interceptor
func StreamClientLog2(logger *zap.Logger, opts ...grpc_zap.Option) grpc.StreamClientInterceptor {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	return grpc_zap.StreamClientInterceptor(logger, opts...)
}

// ---------------------------------- server interceptor ----------------------------------

var ignoreLogMethods = map[string]struct{}{} // ignore printing methods

// LogOption log settings
type LogOption func(*logOptions)

type logOptions struct {
	fields        map[string]interface{}
	ignoreMethods map[string]struct{}
}

func defaultLogOptions() *logOptions {
	return &logOptions{
		fields:        make(map[string]interface{}),
		ignoreMethods: make(map[string]struct{}),
	}
}

func (o *logOptions) apply(opts ...LogOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithLogFields adding a custom print field
func WithLogFields(kvs map[string]interface{}) LogOption {
	return func(o *logOptions) {
		if len(kvs) == 0 {
			return
		}
		o.fields = kvs
	}
}

// WithLogIgnoreMethods ignore printing methods
// fullMethodName format: /packageName.serviceName/methodName,
// example /api.userExample.v1.userExampleService/GetByID
func WithLogIgnoreMethods(fullMethodNames ...string) LogOption {
	return func(o *logOptions) {
		for _, method := range fullMethodNames {
			o.ignoreMethods[method] = struct{}{}
		}
	}
}

// UnaryServerLog server-side log unary interceptor
func UnaryServerLog(logger *zap.Logger, opts ...LogOption) grpc.UnaryServerInterceptor {
	o := defaultLogOptions()
	o.apply(opts...)
	ignoreLogMethods = o.ignoreMethods

	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	grpc_zap.ReplaceGrpcLoggerV2(logger)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// ignore printing of the specified method
		if _, ok := ignoreLogMethods[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		startTime := time.Now()
		requestID := ServerCtxRequestID(ctx)

		fields := []zap.Field{
			zap.String("rpc_type", "unary"),
			zap.String("full_method", info.FullMethod),
			zap.Any("request", req),
		}
		if requestID != "" {
			fields = append(fields, zap.String(ContextRequestIDKey, requestID))
		}
		logger.Info("<<<<", fields...)

		resp, err := handler(ctx, req)

		data, _ := json.Marshal(resp)
		if len(data) > 300 {
			data = append(data[:300], []byte("......")...)
		}

		fields = []zap.Field{
			zap.String("code", status.Code(err).String()),
			zap.Error(err),
			zap.String("rpc_type", "unary"),
			zap.String("full_method", info.FullMethod),
			zap.String("response", string(data)),
			zap.Int64("time_us", time.Since(startTime).Microseconds()),
		}
		if requestID != "" {
			fields = append(fields, zap.String(ContextRequestIDKey, requestID))
		}
		logger.Info(">>>>", fields...)

		return resp, err
	}
}

// UnaryServerLog2 server-side log unary interceptor
func UnaryServerLog2(logger *zap.Logger, opts ...LogOption) grpc.UnaryServerInterceptor {
	o := defaultLogOptions()
	o.apply(opts...)
	ignoreLogMethods = o.ignoreMethods

	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	grpc_zap.ReplaceGrpcLoggerV2(logger)

	// log settings, default printing of client disconnection information, example https://pkg.go.dev/github.com/grpc-ecosystem/go-grpc-middleware/logging/zap
	zapOptions := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_us", duration.Microseconds())
		}),
	}

	// custom log fields
	for key, val := range o.fields {
		zapOptions = append(zapOptions, grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Any(key, val)
		}))
	}

	// custom call method for skipping log
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

// UnaryServerCtxTags extractor field unary interceptor
//func UnaryServerCtxTags() grpc.UnaryServerInterceptor {
//	return grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor))
//}

// StreamServerLog Server-side log stream interceptor
func StreamServerLog(logger *zap.Logger, opts ...LogOption) grpc.StreamServerInterceptor {
	o := defaultLogOptions()
	o.apply(opts...)
	ignoreLogMethods = o.ignoreMethods

	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	grpc_zap.ReplaceGrpcLoggerV2(logger)

	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// ignore printing of the specified method
		if _, ok := ignoreLogMethods[info.FullMethod]; ok {
			return handler(srv, stream)
		}

		startTime := time.Now()
		requestID := ServerCtxRequestID(stream.Context())

		fields := []zap.Field{
			zap.String("rpc_type", "stream"),
			zap.String("full_method", info.FullMethod),
		}
		if requestID != "" {
			fields = append(fields, zap.String(ContextRequestIDKey, requestID))
		}
		logger.Info("<<<<", fields...)

		err := handler(srv, stream)

		fields = []zap.Field{
			zap.String("code", status.Code(err).String()),
			zap.String("rpc_type", "stream"),
			zap.String("full_method", info.FullMethod),
			zap.Int64("time_us", time.Since(startTime).Microseconds()),
		}
		if requestID != "" {
			fields = append(fields, zap.String(ContextRequestIDKey, requestID))
		}
		logger.Info(">>>>", fields...)

		return err
	}
}

// StreamServerLog2 Server-side log stream interceptor
func StreamServerLog2(logger *zap.Logger, opts ...LogOption) grpc.StreamServerInterceptor {
	o := defaultLogOptions()
	o.apply(opts...)
	ignoreLogMethods = o.ignoreMethods

	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	grpc_zap.ReplaceGrpcLoggerV2(logger)

	// log settings, default printing of client disconnection information, example https://pkg.go.dev/github.com/grpc-ecosystem/go-grpc-middleware/logging/zap
	zapOptions := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_us", duration.Microseconds())
		}),
	}

	// custom log fields
	for key, val := range o.fields {
		zapOptions = append(zapOptions, grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Any(key, val)
		}))
	}

	// custom call method for skipping log
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

// StreamServerCtxTags extractor field stream interceptor
//func StreamServerCtxTags() grpc.StreamServerInterceptor {
//	return grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor))
//}
