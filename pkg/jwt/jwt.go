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
	Name string `json:"name"`
	jwt.RegisteredClaims
}

// GenerateToken generate token by uid and name
func GenerateToken(uid string, name ...string) (string, error) {
	if opt == nil {
		return "", errInit
	}

	nameVal := ""
	if len(name) > 0 {
		nameVal = name[0]
	}
	claims := Claims{
		uid,
		nameVal,
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

// RefreshToken refresh token
func RefreshToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(opt.expire))
	claims.RegisteredClaims.IssuedAt = jwt.NewNumericDate(time.Now())
	token := jwt.NewWithClaims(opt.signingMethod, claims)
	return token.SignedString(opt.signingKey)
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

// RefreshCustomToken refresh custom token
func RefreshCustomToken(tokenString string) (string, error) {
	claims, err := ParseCustomToken(tokenString)
	if err != nil {
		return "", err
	}
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(opt.expire))
	claims.RegisteredClaims.IssuedAt = jwt.NewNumericDate(time.Now())
	token := jwt.NewWithClaims(opt.signingMethod, claims)
	return token.SignedString(opt.signingKey)
}
