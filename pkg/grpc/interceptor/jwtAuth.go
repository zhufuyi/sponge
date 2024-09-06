package interceptor

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/zhufuyi/sponge/pkg/jwt"
)

// ---------------------------------- client ----------------------------------

// SetJwtTokenToCtx set the token (excluding prefix Bearer) to the context in grpc client side
// Example:
//
// authorization := "Bearer jwt-token"
//
//	ctx := SetJwtTokenToCtx(ctx, authorization)
//	cli.GetByID(ctx, req)
func SetJwtTokenToCtx(ctx context.Context, token string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		md.Set(headerAuthorize, authScheme+" "+token)
	} else {
		md = metadata.Pairs(headerAuthorize, authScheme+" "+token)
	}
	return metadata.NewOutgoingContext(ctx, md)
}

// SetAuthToCtx set the authorization (including prefix Bearer) to the context in grpc client side
// Example:
//
//	ctx := SetAuthToCtx(ctx, authorization)
//	cli.GetByID(ctx, req)
func SetAuthToCtx(ctx context.Context, authorization string) context.Context {
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

// GetAuthorization combining tokens into authentication information
func GetAuthorization(token string) string {
	return authScheme + " " + token
}

// GetAuthCtxKey get the name of Claims
func GetAuthCtxKey() string {
	return authCtxClaimsName
}

// StandardVerifyFn verify function, tokenTail32 is the last 32 characters of the token.
type StandardVerifyFn = func(claims *jwt.Claims, tokenTail32 string) error

// CustomVerifyFn verify custom function, tokenTail32 is the last 32 characters of the token.
type CustomVerifyFn = func(claims *jwt.CustomClaims, tokenTail32 string) error

type verifyOptions struct {
	verifyType       int // 1: use StandardVerifyFn, 2:use CustomVerifyFn
	standardVerifyFn StandardVerifyFn
	customVerifyFn   CustomVerifyFn
}

func defaultVerifyOptions() *verifyOptions {
	return &verifyOptions{
		verifyType: 1,
	}
}

// AuthOption setting the Authentication Field
type AuthOption func(*authOptions)

// authOptions settings
type authOptions struct {
	authScheme    string
	ctxClaimsName string
	ignoreMethods map[string]struct{}

	verifyOpts *verifyOptions
}

func defaultAuthOptions() *authOptions {
	return &authOptions{
		authScheme:    authScheme,
		ctxClaimsName: authCtxClaimsName,
		ignoreMethods: make(map[string]struct{}), // ways to ignore forensics

		verifyOpts: defaultVerifyOptions(),
	}
}

func (o *authOptions) apply(opts ...AuthOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithAuthScheme set the message prefix for authentication
func WithAuthScheme(scheme string) AuthOption {
	return func(o *authOptions) {
		o.authScheme = scheme
	}
}

// WithAuthClaimsName set the key name of the information in ctx for authentication
func WithAuthClaimsName(claimsName string) AuthOption {
	return func(o *authOptions) {
		o.ctxClaimsName = claimsName
	}
}

// WithAuthIgnoreMethods ways to ignore forensics
// fullMethodName format: /packageName.serviceName/methodName,
// example /api.userExample.v1.userExampleService/GetByID
func WithAuthIgnoreMethods(fullMethodNames ...string) AuthOption {
	return func(o *authOptions) {
		for _, method := range fullMethodNames {
			o.ignoreMethods[method] = struct{}{}
		}
	}
}

// WithStandardVerify set the standard verify function for authentication
func WithStandardVerify(verify StandardVerifyFn) AuthOption {
	return func(o *authOptions) {
		if o.verifyOpts == nil {
			o.verifyOpts = defaultVerifyOptions()
		}
		o.verifyOpts.verifyType = 1
		o.verifyOpts.standardVerifyFn = verify
	}
}

// WithCustomVerify set the custom verify function for authentication
func WithCustomVerify(verify CustomVerifyFn) AuthOption {
	return func(o *authOptions) {
		if o.verifyOpts == nil {
			o.verifyOpts = defaultVerifyOptions()
		}
		o.verifyOpts.verifyType = 2
		o.verifyOpts.customVerifyFn = verify
	}
}

// -------------------------------------------------------------------------------------------

// verify authorization from context, support standard and custom verify processing
func jwtVerify(ctx context.Context, opt *verifyOptions) (context.Context, error) {
	if opt == nil {
		opt = &verifyOptions{
			verifyType: 1, // default use VerifyGeneralFn
		}
	}

	token, err := grpc_auth.AuthFromMD(ctx, authScheme) // key is authScheme
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "%v", err)
	}

	if len(token) <= 100 {
		return ctx, status.Errorf(codes.Unauthenticated, "authorization is illegal")
	}

	// custom claims
	if opt.verifyType == 2 {
		var claims *jwt.CustomClaims
		claims, err = jwt.ParseCustomToken(token)
		if err != nil {
			return ctx, status.Errorf(codes.Unauthenticated, "%v", err)
		}
		if opt.customVerifyFn != nil {
			tokenTail32 := token[len(token)-16:]
			err = opt.customVerifyFn(claims, tokenTail32)
			if err != nil {
				return ctx, status.Errorf(codes.Unauthenticated, "%v", err)
			}
		}

		newCtx := context.WithValue(ctx, authCtxClaimsName, claims) //nolint
		return newCtx, nil
	}

	// standard claims
	claims, err := jwt.ParseToken(token)
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "%v", err)
	}
	if opt.standardVerifyFn != nil {
		tokenTail32 := token[len(token)-16:]
		err = opt.standardVerifyFn(claims, tokenTail32)
		if err != nil {
			return ctx, status.Errorf(codes.Unauthenticated, "%v", err)
		}
	}
	newCtx := context.WithValue(ctx, authCtxClaimsName, claims) //nolint
	return newCtx, nil
}

// GetJwtClaims get the jwt standard claims from context, contains fixed fields uid and name
func GetJwtClaims(ctx context.Context) (*jwt.Claims, bool) {
	v, ok := ctx.Value(authCtxClaimsName).(*jwt.Claims)
	return v, ok
}

// GetJwtCustomClaims get the jwt custom claims from context, contains custom fields
func GetJwtCustomClaims(ctx context.Context) (*jwt.CustomClaims, bool) {
	v, ok := ctx.Value(authCtxClaimsName).(*jwt.CustomClaims)
	return v, ok
}

// UnaryServerJwtAuth jwt unary interceptor
func UnaryServerJwtAuth(opts ...AuthOption) grpc.UnaryServerInterceptor {
	o := defaultAuthOptions()
	o.apply(opts...)
	authScheme = o.authScheme
	authCtxClaimsName = o.ctxClaimsName
	authIgnoreMethods = o.ignoreMethods
	verifyOpt := o.verifyOpts

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		var err error

		if _, ok := authIgnoreMethods[info.FullMethod]; ok {
			newCtx = ctx
		} else {
			newCtx, err = jwtVerify(ctx, verifyOpt)
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
	verifyOpt := o.verifyOpts

	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var newCtx context.Context
		var err error

		if _, ok := authIgnoreMethods[info.FullMethod]; ok {
			newCtx = stream.Context()
		} else {
			newCtx, err = jwtVerify(stream.Context(), verifyOpt)
			if err != nil {
				return err
			}
		}

		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}
