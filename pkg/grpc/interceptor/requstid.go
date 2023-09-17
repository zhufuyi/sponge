package interceptor

import (
	"context"

	"github.com/zhufuyi/sponge/pkg/krand"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	// ContextRequestIDKey context request id for context
	ContextRequestIDKey = "request_id"
)

// CtxKeyString for context.WithValue key type
type CtxKeyString string

// RequestIDKey "request_id"
var RequestIDKey = CtxKeyString(ContextRequestIDKey)

// ---------------------------------- client interceptor ----------------------------------

// ClientCtxRequestID get request id from rpc client context.Context
func ClientCtxRequestID(ctx context.Context) string {
	return metautils.ExtractOutgoing(ctx).Get(ContextRequestIDKey)
}

// ClientCtxRequestIDField get request id field from rpc client context.Context
func ClientCtxRequestIDField(ctx context.Context) zap.Field {
	return zap.String(ContextRequestIDKey, ClientCtxRequestID(ctx))
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

// ServerCtxRequestID get request id from rpc server context.Context
func ServerCtxRequestID(ctx context.Context) string {
	return metautils.ExtractIncoming(ctx).Get(ContextRequestIDKey)
}

// ServerCtxRequestIDField get request id field from rpc server context.Context
func ServerCtxRequestIDField(ctx context.Context) zap.Field {
	return zap.String(ContextRequestIDKey, ServerCtxRequestID(ctx))
}

// UnaryServerRequestID server-side request_id unary interceptor
func UnaryServerRequestID() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		requestID := ServerCtxRequestID(ctx)
		if requestID == "" {
			requestID = krand.String(krand.R_All, 10)
			ctx = metautils.ExtractIncoming(ctx).Add(ContextRequestIDKey, requestID).ToIncoming(ctx)
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
		//	ctx = metautils.ExtractIncoming(ctx).Add(ContextRequestIDKey, requestID).ToIncoming(ctx)
		//}
		return handler(srv, stream)
	}
}
