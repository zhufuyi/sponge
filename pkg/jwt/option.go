package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// HS256 Method
	HS256 = jwt.SigningMethodHS256
	// HS384 Method
	HS384 = jwt.SigningMethodHS384
	// HS512 Method
	HS512 = jwt.SigningMethodHS512
)

var (
	defaultSigningKey    = []byte("zaq12wsxmko0") // default key
	defaultSigningMethod = HS256                  // default HS256
	defaultExpire        = 24 * time.Hour         // default expiration
	defaultIssuer        = ""
)

type options struct {
	signingKey    []byte
	expire        time.Duration
	issuer        string
	signingMethod *jwt.SigningMethodHMAC
}

func defaultOptions() *options {
	return &options{
		signingKey:    defaultSigningKey,
		signingMethod: defaultSigningMethod,
		expire:        defaultExpire,
		issuer:        defaultIssuer,
	}
}

// Option set the jwt options.
type Option func(*options)

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithSigningKey set signing key value
func WithSigningKey(key string) Option {
	return func(o *options) {
		o.signingKey = []byte(key)
	}
}

// WithSigningMethod set signing method value
func WithSigningMethod(sm *jwt.SigningMethodHMAC) Option {
	return func(o *options) {
		o.signingMethod = sm
	}
}

// WithExpire set expire value
func WithExpire(d time.Duration) Option {
	return func(o *options) {
		o.expire = d
	}
}

// WithIssuer set issuer value
func WithIssuer(issuer string) Option {
	return func(o *options) {
		o.issuer = issuer
	}
}

var (
	errSignature = errors.New("signature failure")
	errInit      = errors.New("not yet initialized jwt, usage 'jwt.Init()'")
)
