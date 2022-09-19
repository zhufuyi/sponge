package interceptor

import (
	"context"

	"github.com/zhufuyi/sponge/pkg/jwt"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ---------------------------------- server interceptor ----------------------------------

var (
	// auth Scheme
	authScheme = "Bearer"

	// 鉴权信息在ctx中key名
	authCtxClaimsName = "tokenInfo"

	// 跳过认证方法集合
	authIgnoreMethods = map[string]struct{}{}
)

// AuthOption 鉴权设置
type AuthOption func(*AuthOptions)

type AuthOptions struct {
	authScheme    string
	ctxClaimsName string
	ignoreMethods map[string]struct{}
}

func defaultAuthOptions() *AuthOptions {
	return &AuthOptions{
		authScheme:    authScheme,
		ctxClaimsName: authCtxClaimsName,
		ignoreMethods: make(map[string]struct{}), // 忽略鉴权的方法
	}
}

func (o *AuthOptions) apply(opts ...AuthOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithAuthScheme 设置鉴权的信息前缀
func WithAuthScheme(scheme string) AuthOption {
	return func(o *AuthOptions) {
		o.authScheme = scheme
	}
}

// WithAuthClaimsName 设置鉴权的信息在ctx的key名称
func WithAuthClaimsName(claimsName string) AuthOption {
	return func(o *AuthOptions) {
		o.ctxClaimsName = claimsName
	}
}

// WithAuthIgnoreMethods 忽略鉴权的方法
// fullMethodName格式: /packageName.serviceName/methodName，
// 示例/api.userExample.v1.userExampleService/GetByID
func WithAuthIgnoreMethods(fullMethodNames ...string) AuthOption {
	return func(o *AuthOptions) {
		for _, method := range fullMethodNames {
			o.ignoreMethods[method] = struct{}{}
		}
	}
}

// GetAuthorization 根据token组合成鉴权信息
func GetAuthorization(token string) string {
	return authScheme + " " + token
}

// GetAuthCtxKey 获取Claims的名称
func GetAuthCtxKey() string {
	return authCtxClaimsName
}

// JwtVerify 从context获取authorization来验证是否合法，authorization组成格式：authScheme token
func JwtVerify(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, authScheme)
	if err != nil {
		return nil, err
	}

	cc, err := jwt.VerifyToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%v", err)
	}

	newCtx := context.WithValue(ctx, authCtxClaimsName, cc) //nolint 后面方法可以通过ctx.Value(interceptor.GetAuthCtxKey()).(*jwt.CustomClaims)

	return newCtx, nil
}

// UnaryServerJwtAuth jwt鉴权unary拦截器
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

// StreamServerJwtAuth jwt鉴权stream拦截器
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
