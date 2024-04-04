package interceptor

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	zapLog "github.com/zhufuyi/sponge/pkg/logger"
)

// ---------------------------------- client interceptor ----------------------------------

// UnaryClientLog client log unary interceptor
func UnaryClientLog(logger *zap.Logger, opts ...LogOption) grpc.UnaryClientInterceptor {
	o := defaultLogOptions()
	o.apply(opts...)
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	if o.isReplaceGRPCLogger {
		zapLog.ReplaceGRPCLoggerV2(logger)
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
			zap.String("type", "unary"),
			zap.String("method", method),
			zap.Int64("time_us", time.Since(startTime).Microseconds()),
			reqIDField,
		}

		logger.Info("invoker result", fields...)
		return err
	}
}

// StreamClientLog client log stream interceptor
func StreamClientLog(logger *zap.Logger, opts ...LogOption) grpc.StreamClientInterceptor {
	o := defaultLogOptions()
	o.apply(opts...)
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	if o.isReplaceGRPCLogger {
		zapLog.ReplaceGRPCLoggerV2(logger)
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
			zap.String("type", "stream"),
			zap.String("method", method),
			zap.Int64("time_us", time.Since(startTime).Microseconds()),
			reqIDField,
		}

		logger.Info("invoker result", fields...)
		return clientStream, err
	}
}

// ---------------------------------- server interceptor ----------------------------------

var ignoreLogMethods = map[string]struct{}{} // ignore printing methods

// LogOption log settings
type LogOption func(*logOptions)

type logOptions struct {
	fields              map[string]interface{}
	ignoreMethods       map[string]struct{}
	isReplaceGRPCLogger bool
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

// WithReplaceGRPCLogger replace grpc logger v2
func WithReplaceGRPCLogger() LogOption {
	return func(o *logOptions) {
		o.isReplaceGRPCLogger = true
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
	if o.isReplaceGRPCLogger {
		zapLog.ReplaceGRPCLoggerV2(logger)
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// ignore printing of the specified method
		if _, ok := ignoreLogMethods[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		startTime := time.Now()
		requestID := ServerCtxRequestID(ctx)

		fields := []zap.Field{
			zap.String("type", "unary"),
			zap.String("method", info.FullMethod),
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
			zap.String("type", "unary"),
			zap.String("method", info.FullMethod),
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

// UnaryServerSimpleLog server-side log unary interceptor, only print response
func UnaryServerSimpleLog(logger *zap.Logger, opts ...LogOption) grpc.UnaryServerInterceptor {
	o := defaultLogOptions()
	o.apply(opts...)
	ignoreLogMethods = o.ignoreMethods

	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	if o.isReplaceGRPCLogger {
		zapLog.ReplaceGRPCLoggerV2(logger)
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// ignore printing of the specified method
		if _, ok := ignoreLogMethods[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		startTime := time.Now()
		requestID := ServerCtxRequestID(ctx)

		resp, err := handler(ctx, req)

		fields := []zap.Field{
			zap.String("code", status.Code(err).String()),
			zap.Error(err),
			zap.String("type", "unary"),
			zap.String("method", info.FullMethod),
			zap.Int64("time_us", time.Since(startTime).Microseconds()),
		}
		if requestID != "" {
			fields = append(fields, zap.String(ContextRequestIDKey, requestID))
		}
		logger.Info("[GRPC]", fields...)

		return resp, err
	}
}

// StreamServerLog Server-side log stream interceptor
func StreamServerLog(logger *zap.Logger, opts ...LogOption) grpc.StreamServerInterceptor {
	o := defaultLogOptions()
	o.apply(opts...)
	ignoreLogMethods = o.ignoreMethods

	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	if o.isReplaceGRPCLogger {
		zapLog.ReplaceGRPCLoggerV2(logger)
	}

	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// ignore printing of the specified method
		if _, ok := ignoreLogMethods[info.FullMethod]; ok {
			return handler(srv, stream)
		}

		startTime := time.Now()
		requestID := ServerCtxRequestID(stream.Context())

		fields := []zap.Field{
			zap.String("type", "stream"),
			zap.String("method", info.FullMethod),
		}
		if requestID != "" {
			fields = append(fields, zap.String(ContextRequestIDKey, requestID))
		}
		logger.Info("<<<<", fields...)

		err := handler(srv, stream)

		fields = []zap.Field{
			zap.String("code", status.Code(err).String()),
			zap.String("type", "stream"),
			zap.String("method", info.FullMethod),
			zap.Int64("time_us", time.Since(startTime).Microseconds()),
		}
		if requestID != "" {
			fields = append(fields, zap.String(ContextRequestIDKey, requestID))
		}
		logger.Info(">>>>", fields...)

		return err
	}
}

// StreamServerSimpleLog Server-side log stream interceptor, only print response
func StreamServerSimpleLog(logger *zap.Logger, opts ...LogOption) grpc.StreamServerInterceptor {
	o := defaultLogOptions()
	o.apply(opts...)
	ignoreLogMethods = o.ignoreMethods

	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	if o.isReplaceGRPCLogger {
		zapLog.ReplaceGRPCLoggerV2(logger)
	}

	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// ignore printing of the specified method
		if _, ok := ignoreLogMethods[info.FullMethod]; ok {
			return handler(srv, stream)
		}

		startTime := time.Now()
		requestID := ServerCtxRequestID(stream.Context())

		err := handler(srv, stream)

		fields := []zap.Field{
			zap.String("code", status.Code(err).String()),
			zap.String("type", "stream"),
			zap.String("method", info.FullMethod),
			zap.Int64("time_us", time.Since(startTime).Microseconds()),
		}
		if requestID != "" {
			fields = append(fields, zap.String(ContextRequestIDKey, requestID))
		}
		logger.Info("[GRPC]", fields...)

		return err
	}
}
