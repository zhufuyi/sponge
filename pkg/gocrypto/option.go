package gocrypto

import "crypto"

const (
	modeECB = "ECB"
	modeCBC = "CBC"
	modeCFB = "CFB"
	modeCTR = "CTR"
)

var (
	defaultAesKey = []byte("mKoF_pL,NjI9=I;w") // aes key
	defaultDesKey = []byte("VgY7*uHb")         // des key
	defaultMode   = "ECB"

	defaultRsaFormat   = "PKCS#1"
	defaultRsaHashType = crypto.SHA1
)

type aesOptions struct {
	// the length of the key must be one of 16,24,32, corresponding to
	// AES-128,AES-192,AES-256 respectively.
	aesKey []byte
	// there are four operating modes in total, ECB CBC CFB CTR
	mode string
}

// AesOption set the aes options.
type AesOption func(*aesOptions)

func (o *aesOptions) apply(opts ...AesOption) {
	for _, opt := range opts {
		opt(o)
	}
}

func defaultAesOptions() *aesOptions {
	return &aesOptions{
		aesKey: defaultAesKey,
		mode:   defaultMode,
	}
}

// WithAesKey set aes key
func WithAesKey(key []byte) AesOption {
	return func(o *aesOptions) {
		o.aesKey = key
	}
}

// WithAesModeCBC set mode to CBC
func WithAesModeCBC() AesOption {
	return func(o *aesOptions) {
		o.mode = modeCBC
	}
}

// WithAesModeECB set mode to ECB
func WithAesModeECB() AesOption {
	return func(o *aesOptions) {
		o.mode = modeECB
	}
}

// WithAesModeCFB set mode to CFB
func WithAesModeCFB() AesOption {
	return func(o *aesOptions) {
		o.mode = modeCFB
	}
}

// WithAesModeCTR set mode to CTR
func WithAesModeCTR() AesOption {
	return func(o *aesOptions) {
		o.mode = modeCTR
	}
}

// ------------------------------------------------------------------------------------------

type desOptions struct {
	desKey []byte // the length of the key must be 8
	mode   string // there are four operating modes in total, ECB CBC CFB CTR
}

// DesOption set the des options.
type DesOption func(*desOptions)

func (o *desOptions) apply(opts ...DesOption) {
	for _, opt := range opts {
		opt(o)
	}
}

func defaultDesOptions() *desOptions {
	return &desOptions{
		desKey: defaultDesKey,
		mode:   defaultMode,
	}
}

// WithDesKey set des key
func WithDesKey(key []byte) DesOption {
	return func(o *desOptions) {
		o.desKey = key
	}
}

// WithDesModeCBC set mode to CBC
func WithDesModeCBC() DesOption {
	return func(o *desOptions) {
		o.mode = modeCBC
	}
}

// WithDesModeECB set mode to ECB
func WithDesModeECB() DesOption {
	return func(o *desOptions) {
		o.mode = modeECB
	}
}

// WithDesModeCFB set mode to CFB
func WithDesModeCFB() DesOption {
	return func(o *desOptions) {
		o.mode = modeCFB
	}
}

// WithDesModeCTR set mode to CTR
func WithDesModeCTR() DesOption {
	return func(o *desOptions) {
		o.mode = modeCTR
	}
}

// ------------------------------------------------------------------------------------------

type rsaOptions struct {
	// rsa key pair format
	format string
	// hash types for signatures and signature verification
	hashType crypto.Hash
}

// RsaOption set the rsa options.
type RsaOption func(*rsaOptions)

func (o *rsaOptions) apply(opts ...RsaOption) {
	for _, opt := range opts {
		opt(o)
	}
}

func defaultRsaOptions() *rsaOptions {
	return &rsaOptions{
		format:   defaultRsaFormat,
		hashType: defaultRsaHashType,
	}
}

// WithRsaFormatPKCS1 set format
func WithRsaFormatPKCS1() RsaOption {
	return func(o *rsaOptions) {
		o.format = pkcs1
	}
}

// WithRsaFormatPKCS8 set format
func WithRsaFormatPKCS8() RsaOption {
	return func(o *rsaOptions) {
		o.format = pkcs8
	}
}

// WithRsaHashTypeMd5 set hash type
func WithRsaHashTypeMd5() RsaOption {
	return func(o *rsaOptions) {
		o.hashType = crypto.MD5
	}
}

// WithRsaHashTypeSha1 set hash type
func WithRsaHashTypeSha1() RsaOption {
	return func(o *rsaOptions) {
		o.hashType = crypto.SHA1
	}
}

// WithRsaHashTypeSha256 set hash type
func WithRsaHashTypeSha256() RsaOption {
	return func(o *rsaOptions) {
		o.hashType = crypto.SHA256
	}
}

// WithRsaHashTypeSha512 set hash type
func WithRsaHashTypeSha512() RsaOption {
	return func(o *rsaOptions) {
		o.hashType = crypto.SHA512
	}
}

// WithRsaHashType set hash type
func WithRsaHashType(hash crypto.Hash) RsaOption {
	return func(o *rsaOptions) {
		o.hashType = hash
	}
}
