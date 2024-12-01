package interceptor

import (
	"context"
	"testing"

	grpc_metadata "github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestToken(t *testing.T) {
	to := &authToken{
		AppID:    "grpc",
		AppKey:   "123456",
		IsSecure: false,
	}

	metadata, err := to.GetRequestMetadata(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, metadata)

	isSecure := to.RequireTransportSecurity()
	assert.Equal(t, to.IsSecure, isSecure)
}

func TestClientTokenOption(t *testing.T) {
	option := ClientTokenOption("grpc", "123456", false)
	assert.NotNil(t, option)

}

func TestUnaryServerToken(t *testing.T) {
	f := func(appID string, appKey string) error {
		if appID != "grpc" || appKey != "123456" {
			return status.Errorf(codes.Unauthenticated, "app id or app key checksum failure")
		}
		return nil
	}
	interceptor := UnaryServerToken(f)
	assert.NotNil(t, interceptor)

	ctx := context.Background()
	_, err := interceptor(ctx, nil, unaryServerInfo, unaryServerHandler)
	assert.NotNil(t, err)
	ctx = grpc_metadata.ExtractIncoming(ctx).Add("app_id", "grpc").ToIncoming(ctx)
	ctx = grpc_metadata.ExtractIncoming(ctx).Add("app_key", "123456").ToIncoming(ctx)
	_, err = interceptor(ctx, nil, unaryServerInfo, unaryServerHandler)
	assert.NoError(t, err)
}

func TestStreamServerToken(t *testing.T) {
	f := func(appID string, appKey string) error {
		if appID != "grpc" || appKey != "123456" {
			return status.Errorf(codes.Unauthenticated, "app id or app key checksum failure")
		}
		return nil
	}
	interceptor := StreamServerToken(f)
	assert.NotNil(t, interceptor)

	ctx := context.Background()
	err := interceptor(nil, newStreamServer(ctx), streamServerInfo, streamServerHandler)
	assert.NotNil(t, err)
	ctx = grpc_metadata.ExtractIncoming(ctx).Add("app_id", "grpc").ToIncoming(ctx)
	ctx = grpc_metadata.ExtractIncoming(ctx).Add("app_key", "123456").ToIncoming(ctx)
	err = interceptor(nil, newStreamServer(ctx), streamServerInfo, streamServerHandler)
	assert.NoError(t, err)
}
