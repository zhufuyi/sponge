package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(WithSigningKey("foo"))
}

func TestWithExpire(t *testing.T) {
	testData := time.Second * 3
	opt := WithExpire(testData)
	o := new(options)
	o.apply(opt)
	assert.Equal(t, testData, o.expire)
}

func TestWithIssuer(t *testing.T) {
	testData := "issuer"
	opt := WithIssuer(testData)
	o := new(options)
	o.apply(opt)
	assert.Equal(t, testData, o.issuer)
}

func TestWithSigningKey(t *testing.T) {
	testData := "key"
	opt := WithSigningKey(testData)
	o := new(options)
	o.apply(opt)
	assert.Equal(t, testData, string(o.signingKey))
}

func TestWithSigningMethod(t *testing.T) {
	testData := jwt.SigningMethodHS384
	opt := WithSigningMethod(testData)
	o := new(options)
	o.apply(opt)
	assert.Equal(t, testData, o.signingMethod)
}

func Test_defaultOptions(t *testing.T) {
	o := defaultOptions()
	assert.NotNil(t, o)
}

func Test_options_apply(t *testing.T) {
	testData := "key"
	opt := WithSigningKey(testData)
	o := new(options)
	o.apply(opt)
	assert.Equal(t, testData, string(o.signingKey))
}
