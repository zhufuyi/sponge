package interceptor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAuthCtxKey(t *testing.T) {
	key := GetAuthCtxKey()
	assert.Equal(t, authCtxClaimsName, key)
}

func TestGetAuthorization(t *testing.T) {
	testData := "token"
	authorization := GetAuthorization(testData)
	assert.Equal(t, authScheme+" "+testData, authorization)
}

func TestJwtVerify(t *testing.T) {
	ctx := context.WithValue(context.Background(), "authorization", authScheme+" eyJhbGciOi......5cCI6Ikp")
	_, err := JwtVerify(ctx)
	assert.NotNil(t, err)
}

func TestStreamServerJwtAuth(t *testing.T) {
	interceptor := StreamServerJwtAuth()
	assert.NotNil(t, interceptor)
}

func TestUnaryServerJwtAuth(t *testing.T) {
	interceptor := UnaryServerJwtAuth()
	assert.NotNil(t, interceptor)
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
