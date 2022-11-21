package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// CustomClaims add custom fields to StandardClaims' payload
type CustomClaims struct {
	UID  string `json:"uid"`
	Role string `json:"role"`
	jwt.StandardClaims
}

// GenerateToken generate token
func GenerateToken(uid string, role ...string) (string, error) {
	if opt == nil {
		return "", errInit
	}

	roleVal := ""
	if len(role) > 0 {
		roleVal = role[0]
	}
	claims := CustomClaims{
		uid,
		roleVal,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(opt.expire).Unix(),
			Issuer:    opt.issuer,
		},
	}

	token := jwt.NewWithClaims(opt.signingMethod, claims)
	return token.SignedString(opt.signingKey)
}

// VerifyToken verify token
func VerifyToken(tokenString string) (*CustomClaims, error) {
	if opt == nil {
		return nil, errInit
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return opt.signingKey, nil
	})
	if err != nil {
		ve, ok := err.(*jwt.ValidationError)
		if ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errFormat
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errExpired
			} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
				return nil, errUnverifiable
			} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				return nil, errSignature
			} else {
				return nil, ve
			}
		}
		return nil, errSignature
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errSignature
}
