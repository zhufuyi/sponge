package interceptor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/go-dev-frame/sponge/pkg/jwt"
	"github.com/go-dev-frame/sponge/pkg/utils"
)

var (
	expectedUid  = "100"
	expectedName = "tom"

	expectedFields = jwt.KV{"id": utils.StrToUint64(expectedUid), "name": expectedName, "age": 10}
)

func standardVerifyHandler(claims *jwt.Claims, tokenTail32 string) error {
	// token := getToken(claims.UID)
	// if  token[len(token)-32:] != tokenTail32 { return err }

	if claims.UID != expectedUid || claims.Name != expectedName {
		return status.Error(codes.Unauthenticated, "id or name not match")
	}

	return nil
}

func customVerifyHandler(claims *jwt.CustomClaims, tokenTail32 string) error {
	err := status.Error(codes.Unauthenticated, "custom verify failed")

	//token, fields := getToken(id)
	// if  token[len(token)-32:] != tokenTail32 { return err }

	id, exist := claims.GetUint64("id")
	if !exist || id != expectedFields["id"] {
		return err
	}

	name, exist := claims.GetString("name")
	if !exist || name != expectedFields["name"] {
		return err
	}

	age, exist := claims.GetInt("age")
	if !exist || age != expectedFields["age"] {
		return err
	}

	return nil
}

func TestJwtVerify(t *testing.T) {
	jwt.Init()
	ctx := context.Background()
	token, _ := jwt.GenerateToken(expectedUid, expectedName)

	// success test
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{headerAuthorize: []string{GetAuthorization(token)}})
	newCtx, err := jwtVerify(ctx, nil)
	assert.NoError(t, err)
	claims, ok := GetJwtClaims(newCtx)
	assert.True(t, ok)
	assert.Equal(t, expectedUid, claims.UID)

	// success test
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{headerAuthorize: []string{GetAuthorization(token)}})
	newCtx, err = jwtVerify(ctx, &verifyOptions{verifyType: 1, standardVerifyFn: standardVerifyHandler})
	assert.NoError(t, err)
	claims, ok = GetJwtClaims(newCtx)
	assert.True(t, ok)
	assert.Equal(t, expectedUid, claims.UID)

	authorization := []string{GetAuthorization("error token......")}
	// authorization format error, missing token
	ctx = metadata.NewIncomingContext(context.Background(), metadata.MD{headerAuthorize: authorization})
	_, err = jwtVerify(ctx, nil)
	assert.Error(t, err)

	// authorization format error, missing Bearer
	ctx = context.WithValue(context.Background(), headerAuthorize, authorization)
	_, err = jwtVerify(ctx, nil)
	assert.Error(t, err)
}

func TestJwtCustomVerify(t *testing.T) {
	jwt.Init()
	ctx := context.Background()
	token, _ := jwt.GenerateCustomToken(expectedFields)
	verifyOpt := &verifyOptions{verifyType: 2}

	// success test
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{headerAuthorize: []string{GetAuthorization(token)}})
	newCtx, err := jwtVerify(ctx, verifyOpt)
	assert.NoError(t, err)
	claims, ok := GetJwtCustomClaims(newCtx)
	assert.True(t, ok)
	assert.Equal(t, expectedName, claims.Fields["name"])

	// success test
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{headerAuthorize: []string{GetAuthorization(token)}})
	verifyOpt.customVerifyFn = customVerifyHandler
	newCtx, err = jwtVerify(ctx, verifyOpt)
	assert.NoError(t, err)
	claims, ok = GetJwtCustomClaims(newCtx)
	assert.True(t, ok)
	assert.Equal(t, expectedName, claims.Fields["name"])

	authorization := []string{GetAuthorization("mock token......")}

	// authorization format error, missing token
	ctx = metadata.NewIncomingContext(context.Background(), metadata.MD{headerAuthorize: authorization})
	_, err = jwtVerify(ctx, verifyOpt)
	assert.Error(t, err)

	// authorization format error, missing Bearer
	ctx = context.WithValue(context.Background(), headerAuthorize, authorization)
	_, err = jwtVerify(ctx, verifyOpt)
	assert.Error(t, err)
}

func TestUnaryServerJwtAuth(t *testing.T) {
	interceptor := UnaryServerJwtAuth(WithStandardVerify(standardVerifyHandler))
	assert.NotNil(t, interceptor)

	// mock client ctx
	jwt.Init()
	token, _ := jwt.GenerateToken(expectedUid, expectedName)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{headerAuthorize: []string{GetAuthorization(token)}})

	_, err := interceptor(ctx, nil, unaryServerInfo, unaryServerHandler)
	assert.NoError(t, err)

	ctx = metadata.NewIncomingContext(context.Background(), metadata.MD{headerAuthorize: []string{GetAuthorization("error token......")}})
	_, err = interceptor(ctx, nil, unaryServerInfo, unaryServerHandler)
	assert.Error(t, err)
}

func TestUnaryServerJwtCustomAuth(t *testing.T) {
	interceptor := UnaryServerJwtAuth(WithCustomVerify(customVerifyHandler))
	assert.NotNil(t, interceptor)

	// mock client ctx
	jwt.Init()
	token, _ := jwt.GenerateCustomToken(expectedFields)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{headerAuthorize: []string{GetAuthorization(token)}})

	_, err := interceptor(ctx, nil, unaryServerInfo, unaryServerHandler)
	assert.NoError(t, err)

	ctx = metadata.NewIncomingContext(context.Background(), metadata.MD{headerAuthorize: []string{GetAuthorization("error token......")}})
	_, err = interceptor(ctx, nil, unaryServerInfo, unaryServerHandler)
	assert.Error(t, err)
}

func TestStreamServerJwtAuth(t *testing.T) {
	interceptor := StreamServerJwtAuth()
	assert.NotNil(t, interceptor)

	jwt.Init()
	token, _ := jwt.GenerateToken(expectedUid, expectedName)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{headerAuthorize: []string{authScheme + " " + token}})
	err := interceptor(nil, newStreamServer(ctx), streamServerInfo, streamServerHandler)
	assert.NoError(t, err)

	err = interceptor(nil, newStreamServer(context.Background()), streamServerInfo, streamServerHandler)
	assert.Error(t, err)
}

func TestStreamServerJwtCustomAuth(t *testing.T) {
	interceptor := StreamServerJwtAuth(WithCustomVerify(nil))
	assert.NotNil(t, interceptor)

	jwt.Init()
	token, _ := jwt.GenerateCustomToken(expectedFields)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{headerAuthorize: []string{authScheme + " " + token}})
	err := interceptor(nil, newStreamServer(ctx), streamServerInfo, streamServerHandler)
	assert.NoError(t, err)

	err = interceptor(nil, newStreamServer(context.Background()), streamServerInfo, streamServerHandler)
	assert.Error(t, err)
}

func TestGetAuthCtxKey(t *testing.T) {
	key := GetAuthCtxKey()
	assert.Equal(t, authCtxClaimsName, key)
}

func TestGetAuthorization(t *testing.T) {
	testData := "token"
	authorization := GetAuthorization(testData)
	assert.Equal(t, authScheme+" "+testData, authorization)
}

func TestAuthOptions(t *testing.T) {
	o := defaultAuthOptions()

	o.apply(WithAuthScheme(authScheme))
	assert.Equal(t, authScheme, o.authScheme)

	o.apply(WithAuthClaimsName(authCtxClaimsName))
	assert.Equal(t, authCtxClaimsName, o.ctxClaimsName)

	o.apply(WithAuthIgnoreMethods("/metrics"))
	assert.Equal(t, struct{}{}, o.ignoreMethods["/metrics"])

	o.apply(WithStandardVerify(nil))
	assert.Equal(t, 1, o.verifyOpts.verifyType)

	o.apply(WithStandardVerify(standardVerifyHandler))
	assert.Equal(t, 1, o.verifyOpts.verifyType)

	o.apply(WithCustomVerify(nil))
	assert.Equal(t, 2, o.verifyOpts.verifyType)

	o.apply(WithCustomVerify(customVerifyHandler))
	assert.Equal(t, 2, o.verifyOpts.verifyType)
}

func TestSetJWTTokenToCtx(t *testing.T) {
	jwt.Init()
	ctx := context.Background()
	token, _ := jwt.GenerateToken(expectedUid, expectedName)
	expected := []string{GetAuthorization(token)}

	ctx = SetJwtTokenToCtx(ctx, token)
	md, _ := metadata.FromOutgoingContext(ctx)
	assert.Equal(t, expected, md.Get(headerAuthorize))
}

func TestSetAuthToCtx(t *testing.T) {
	jwt.Init()
	ctx := context.Background()
	token, _ := jwt.GenerateToken(expectedUid, expectedName)
	authorization := GetAuthorization(token)
	expected := []string{authorization}

	ctx = SetAuthToCtx(ctx, authorization)
	md, _ := metadata.FromOutgoingContext(ctx)
	assert.Equal(t, expected, md.Get(headerAuthorize))
}
