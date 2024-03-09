package interceptor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"

	"github.com/zhufuyi/sponge/pkg/jwt"
)

func TestJwtVerify(t *testing.T) {
	jwt.Init()
	token, _ := jwt.GenerateToken("100")
	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{"authorization": []string{authScheme + " " + token}})
	_, err := JwtVerify(ctx)
	assert.NoError(t, err)

	// token error
	ctx = metadata.NewIncomingContext(context.Background(), metadata.MD{"authorization": []string{authScheme + " " + "token......"}})
	_, err = JwtVerify(ctx)
	assert.Error(t, err)

	// error test
	ctx = context.WithValue(context.Background(), "authorization", "token....")
	_, err = JwtVerify(ctx)
	assert.Error(t, err)
}

func TestUnaryServerJwtAuth(t *testing.T) {
	interceptor := UnaryServerJwtAuth()
	assert.NotNil(t, interceptor)

	jwt.Init()
	token, _ := jwt.GenerateToken("100")
	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{"authorization": []string{authScheme + " " + token}})
	_, err := interceptor(ctx, nil, unaryServerInfo, unaryServerHandler)
	assert.NoError(t, err)

	_, err = interceptor(context.Background(), nil, unaryServerInfo, unaryServerHandler)
	assert.Error(t, err)
}

func TestStreamServerJwtAuth(t *testing.T) {
	interceptor := StreamServerJwtAuth()
	assert.NotNil(t, interceptor)

	jwt.Init()
	token, _ := jwt.GenerateToken("100")
	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{"authorization": []string{authScheme + " " + token}})
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

func TestWithAuthClaimsName(t *testing.T) {
	testData := "demo"
	opt := WithAuthClaimsName(testData)
	o := new(AuthOptions)
	o.apply(opt)
	assert.Equal(t, testData, o.ctxClaimsName)
}

func TestWithAuthIgnoreMethods(t *testing.T) {
	testData := "/method"
	opt := WithAuthIgnoreMethods(testData)
	o := &AuthOptions{ignoreMethods: make(map[string]struct{})}
	o.apply(opt)
	assert.Equal(t, struct{}{}, o.ignoreMethods[testData])
}

func TestWithAuthScheme(t *testing.T) {
	testData := "demo"
	opt := WithAuthScheme(testData)
	o := new(AuthOptions)
	o.apply(opt)
	assert.Equal(t, testData, o.authScheme)
}

func Test_defaultAuthOptions(t *testing.T) {
	o := defaultAuthOptions()
	assert.NotNil(t, o)
}

func TestSetJWTTokenToCtx(t *testing.T) {
	ctx := context.Background()
	expected := []string{"Bearer jwt-token-1"}

	ctx = SetJwtTokenToCtx(ctx, expected[0])
	md, _ := metadata.FromOutgoingContext(ctx)
	assert.Equal(t, expected, md.Get(headerAuthorize))

	expected[0] = "Bearer jwt-token-2"
	ctx = SetJwtTokenToCtx(ctx, expected[0])
	md, _ = metadata.FromOutgoingContext(ctx)
	assert.Equal(t, expected, md.Get(headerAuthorize))
}
