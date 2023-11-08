// Package jwt is token generation and validation.
package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrTokenExpired expired
var ErrTokenExpired = jwt.ErrTokenExpired

var opt *options

// Init initialize jwt
func Init(opts ...Option) {
	o := defaultOptions()
	o.apply(opts...)
	opt = o
}

// Claims my custom claims
type Claims struct {
	UID  string `json:"uid"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken generate token by uid and role
func GenerateToken(uid string, role ...string) (string, error) {
	if opt == nil {
		return "", errInit
	}

	roleVal := ""
	if len(role) > 0 {
		roleVal = role[0]
	}
	claims := Claims{
		uid,
		roleVal,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(opt.expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    opt.issuer,
		},
	}

	token := jwt.NewWithClaims(opt.signingMethod, claims)
	return token.SignedString(opt.signingKey)
}

// ParseToken parse token
func ParseToken(tokenString string) (*Claims, error) {
	if opt == nil {
		return nil, errInit
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return opt.signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errSignature
}

// -------------------------------------------------------------------------------------------

// KV map type
type KV = map[string]interface{}

// CustomClaims custom fields claims
type CustomClaims struct {
	Fields KV `json:"fields"`
	jwt.RegisteredClaims
}

// Get custom field value by key, if not found, return false
func (c *CustomClaims) Get(key string) (val interface{}, isExist bool) {
	if c.Fields == nil {
		return nil, false
	}
	val, isExist = c.Fields[key]
	return val, isExist
}

// GenerateCustomToken generate token by custom fields
func GenerateCustomToken(kv map[string]interface{}) (string, error) {
	if opt == nil {
		return "", errInit
	}

	claims := CustomClaims{
		kv,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(opt.expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    opt.issuer,
		},
	}

	token := jwt.NewWithClaims(opt.signingMethod, claims)
	return token.SignedString(opt.signingKey)
}

// ParseCustomToken parse token
func ParseCustomToken(tokenString string) (*CustomClaims, error) {
	if opt == nil {
		return nil, errInit
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return opt.signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errSignature
}
