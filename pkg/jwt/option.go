package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	// SigningMethodHS256 Method
	SigningMethodHS256 = jwt.SigningMethodHS256
	// SigningMethodHS384 Method
	SigningMethodHS384 = jwt.SigningMethodHS384
	// SigningMethodHS512 Method
	SigningMethodHS512 = jwt.SigningMethodHS512
)

var opt *options

// Init initialize jwt
func Init(opts ...Option) {
	o := defaultOptions()
	o.apply(opts...)
	opt = o
}

var (
	defaultSigningKey    = []byte("zaq12wsxmko0") // default key
	defaultSigningMethod = SigningMethodHS256     // default HS256
	defaultExpire        = 2 * time.Hour          // default expiration
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
	// HS256 Method
	HS256 = jwt.SigningMethodHS256
	// HS384 Method
	HS384 = jwt.SigningMethodHS384
	// HS512 Method
	HS512 = jwt.SigningMethodHS512
)

var (
	errFormat       = errors.New("invalid token format")
	errExpired      = errors.New("token has expired")
	errUnverifiable = errors.New("the token could not be verified due to a signing problem")
	errSignature    = errors.New("signature failure")
	errInit         = errors.New("not yet initialized jwt, usage 'jwt.Init()'")
)
