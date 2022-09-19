package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

/* Standard Claims 的payload没有附加字段 */

// GenerateTokenStandard 生成token
func GenerateTokenStandard() (string, error) {
	if opt == nil {
		return "", errInit
	}

	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(opt.expire).Unix(),
		Issuer:    opt.issuer,
	}

	token := jwt.NewWithClaims(opt.signingMethod, claims)
	return token.SignedString(opt.signingKey)
}

// VerifyTokenStandard 验证token
func VerifyTokenStandard(tokenString string) error {
	if opt == nil {
		return errInit
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return opt.signingKey, nil
	})

	// token有效
	if token.Valid {
		return nil
	}

	ve, ok := err.(*jwt.ValidationError)
	if ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return errFormat
		} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
			return errExpired
		} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
			return errUnverifiable
		} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
			return errSignature
		} else {
			return ve // 其他错误
		}
	}

	return errSignature
}
