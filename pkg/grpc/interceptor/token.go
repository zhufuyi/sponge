package interceptor

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc"
)

// ---------------------------------- client option ----------------------------------

type authToken struct {
	AppID    string `json:"app_id"`
	AppKey   string `json:"app_key"`
	IsSecure bool   `json:"isSecure"`
}

// GetRequestMetadata get metadata
func (t *authToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"app_id":  t.AppID,
		"app_key": t.AppKey,
	}, nil
}

// RequireTransportSecurity is require transport secure
func (t *authToken) RequireTransportSecurity() bool {
	return t.IsSecure
}

// ClientTokenOption client token
func ClientTokenOption(appID string, appKey string, isSecure bool) grpc.DialOption {
	return grpc.WithPerRPCCredentials(&authToken{appID, appKey, isSecure})
}

// ---------------------------------- server interceptor ----------------------------------

// CheckToken check app id and app key
// Example:
//
//	var f CheckToken=func(appID string, appKey string) error{
//		if appID != targetAppID || appKey != targetAppKey {
//			return status.Errorf(codes.Unauthenticated, "app id or app key checksum failure")
//		}
//		return nil
//	}
type CheckToken func(appID string, appKey string) error

// UnaryServerToken recovery unary token
func UnaryServerToken(f CheckToken) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		appID := metautils.ExtractIncoming(ctx).Get("app_id")
		appKey := metautils.ExtractIncoming(ctx).Get("app_key")
		err := f(appID, appKey)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// StreamServerToken recovery stream token
func StreamServerToken(f CheckToken) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		appID := metautils.ExtractIncoming(ctx).Get("app_id")
		appKey := metautils.ExtractIncoming(ctx).Get("app_key")
		err := f(appID, appKey)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}
