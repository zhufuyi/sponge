package interceptor

import (
	"context"
	"sync"

	grpc_metadata "github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zhufuyi/sponge/pkg/krand"
)

var (
	// ContextRequestIDKey request id key for context
	ContextRequestIDKey = "request_id"
	once                sync.Once
)

// SetContextRequestIDKey set context request id key
func SetContextRequestIDKey(key string) {
	if len(key) < 4 {
		return
	}
	once.Do(func() {
		ContextRequestIDKey = key
	})
}

// CtxKeyString for context.WithValue key type
type CtxKeyString string

// RequestIDKey request_id
var RequestIDKey = CtxKeyString(ContextRequestIDKey)

// ---------------------------------- client interceptor ----------------------------------

// CtxRequestIDField get request id field from context.Context
func CtxRequestIDField(ctx context.Context) zap.Field {
	return zap.String(ContextRequestIDKey, grpc_metadata.ExtractOutgoing(ctx).Get(ContextRequestIDKey))
}

// ClientCtxRequestID get request id from rpc client context.Context
func ClientCtxRequestID(ctx context.Context) string {
	return grpc_metadata.ExtractOutgoing(ctx).Get(ContextRequestIDKey)
}

// ClientCtxRequestIDField get request id field from rpc client context.Context
func ClientCtxRequestIDField(ctx context.Context) zap.Field {
	return zap.String(ContextRequestIDKey, grpc_metadata.ExtractOutgoing(ctx).Get(ContextRequestIDKey))
}

// UnaryClientRequestID client-side request_id unary interceptor
func UnaryClientRequestID() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		requestID := ClientCtxRequestID(ctx)
		if requestID == "" {
			requestID = krand.String(krand.R_All, 10)
			ctx = metadata.AppendToOutgoingContext(ctx, ContextRequestIDKey, requestID)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// StreamClientRequestID client request id stream interceptor
func StreamClientRequestID() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
		streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		requestID := ClientCtxRequestID(ctx)
		if requestID == "" {
			requestID = krand.String(krand.R_All, 10)
			ctx = metadata.AppendToOutgoingContext(ctx, ContextRequestIDKey, requestID)
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}

// ---------------------------------- server interceptor ----------------------------------

// KV key value
type KV struct {
	Key string
	Val interface{}
}

// WrapServerCtx wrap context, used in grpc server-side
func WrapServerCtx(ctx context.Context, kvs ...KV) context.Context {
	ctx = context.WithValue(ctx, ContextRequestIDKey, grpc_metadata.ExtractIncoming(ctx).Get(ContextRequestIDKey)) //nolint
	for _, kv := range kvs {
		ctx = context.WithValue(ctx, kv.Key, kv.Val) //nolint
	}
	return ctx
}

// ServerCtxRequestID get request id from rpc server context.Context
func ServerCtxRequestID(ctx context.Context) string {
	return grpc_metadata.ExtractIncoming(ctx).Get(ContextRequestIDKey)
}

// ServerCtxRequestIDField get request id field from rpc server context.Context
func ServerCtxRequestIDField(ctx context.Context) zap.Field {
	return zap.String(ContextRequestIDKey, grpc_metadata.ExtractIncoming(ctx).Get(ContextRequestIDKey))
}

// UnaryServerRequestID server-side request_id unary interceptor
func UnaryServerRequestID() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		requestID := ServerCtxRequestID(ctx)
		if requestID == "" {
			requestID = krand.String(krand.R_All, 10)
			ctx = grpc_metadata.ExtractIncoming(ctx).Add(ContextRequestIDKey, requestID).ToIncoming(ctx)
		}

		return handler(ctx, req)
	}
}

// StreamServerRequestID server-side request id stream interceptor
func StreamServerRequestID() grpc.StreamServerInterceptor {
	// todo
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		//ctx := stream.Context()
		//requestID := ServerCtxRequestID(ctx)
		//if requestID == "" {
		//	requestID = krand.String(krand.R_All, 10)
		//	ctx = grpc_metadata.ExtractIncoming(ctx).Add(ContextRequestIDKey, requestID).ToIncoming(ctx)
		//}
		return handler(srv, stream)
	}
}
