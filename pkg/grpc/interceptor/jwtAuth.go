package interceptor

import (
	"context"

	"github.com/zhufuyi/sponge/pkg/jwt"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ---------------------------------- client ----------------------------------

// SetJwtTokenToCtx set the token to the context in rpc client side
// Example:
//
//	ctx := SetJwtTokenToCtx(ctx, "Bearer jwt-token")
//	cli.GetByID(ctx, req)
func SetJwtTokenToCtx(ctx context.Context, authorization string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		md.Set(headerAuthorize, authorization)
	} else {
		md = metadata.Pairs(headerAuthorize, authorization)
	}
	return metadata.NewOutgoingContext(ctx, md)
}

// ---------------------------------- server interceptor ----------------------------------

var (
	headerAuthorize = "authorization"

	// auth Scheme
	authScheme = "Bearer"

	// authentication information in ctx key name
	authCtxClaimsName = "tokenInfo"

	// collection of skip authentication methods
	authIgnoreMethods = map[string]struct{}{}
)

// AuthOption setting the Authentication Field
type AuthOption func(*AuthOptions)

// AuthOptions settings
type AuthOptions struct {
	authScheme    string
	ctxClaimsName string
	ignoreMethods map[string]struct{}
}

func defaultAuthOptions() *AuthOptions {
	return &AuthOptions{
		authScheme:    authScheme,
		ctxClaimsName: authCtxClaimsName,
		ignoreMethods: make(map[string]struct{}), // ways to ignore forensics
	}
}

func (o *AuthOptions) apply(opts ...AuthOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithAuthScheme set the message prefix for authentication
func WithAuthScheme(scheme string) AuthOption {
	return func(o *AuthOptions) {
		o.authScheme = scheme
	}
}

// WithAuthClaimsName set the key name of the information in ctx for authentication
func WithAuthClaimsName(claimsName string) AuthOption {
	return func(o *AuthOptions) {
		o.ctxClaimsName = claimsName
	}
}

// WithAuthIgnoreMethods ways to ignore forensics
// fullMethodName format: /packageName.serviceName/methodName,
// example /api.userExample.v1.userExampleService/GetByID
func WithAuthIgnoreMethods(fullMethodNames ...string) AuthOption {
	return func(o *AuthOptions) {
		for _, method := range fullMethodNames {
			o.ignoreMethods[method] = struct{}{}
		}
	}
}

// GetAuthorization combining tokens into authentication information
func GetAuthorization(token string) string {
	return authScheme + " " + token
}

// GetAuthCtxKey get the name of Claims
func GetAuthCtxKey() string {
	return authCtxClaimsName
}

// JwtVerify get authorization from context to verify legitimacy, authorization composition format: authScheme token
func JwtVerify(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, authScheme)
	if err != nil {
		return nil, err
	}

	cc, err := jwt.VerifyToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%v", err)
	}

	newCtx := context.WithValue(ctx, authCtxClaimsName, cc) // get value by ctx.Value(interceptor.GetAuthCtxKey()).(*jwt.CustomClaims)

	return newCtx, nil
}

// UnaryServerJwtAuth jwt unary interceptor
func UnaryServerJwtAuth(opts ...AuthOption) grpc.UnaryServerInterceptor {
	o := defaultAuthOptions()
	o.apply(opts...)
	authScheme = o.authScheme
	authCtxClaimsName = o.ctxClaimsName
	authIgnoreMethods = o.ignoreMethods

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		var err error

		if _, ok := authIgnoreMethods[info.FullMethod]; ok {
			newCtx = ctx
		} else {
			newCtx, err = JwtVerify(ctx)
			if err != nil {
				return nil, err
			}
		}

		return handler(newCtx, req)
	}
}

// StreamServerJwtAuth jwt stream interceptor
func StreamServerJwtAuth(opts ...AuthOption) grpc.StreamServerInterceptor {
	o := defaultAuthOptions()
	o.apply(opts...)
	authScheme = o.authScheme
	authCtxClaimsName = o.ctxClaimsName
	authIgnoreMethods = o.ignoreMethods

	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var newCtx context.Context
		var err error

		if _, ok := authIgnoreMethods[info.FullMethod]; ok {
			newCtx = stream.Context()
		} else {
			newCtx, err = JwtVerify(stream.Context())
			if err != nil {
				return err
			}
		}

		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}
