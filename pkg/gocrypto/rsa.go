// Asymmetric encryption and decryption rsa
// 1. public key encryption, private key decryption to get the original text
// 2. private key signature, public key signature verification

package gocrypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
)

const (
	pkcs1 = "PKCS#1"
	pkcs8 = "PKCS#8"
)

// RsaEncrypt rsa encryption, the returned ciphertext is not transcoded
func RsaEncrypt(publicKey []byte, rawData []byte, opts ...RsaOption) ([]byte, error) {
	o := defaultRsaOptions()
	o.apply(opts...)

	return rsaEncryptWithPublicKey(publicKey, rawData)
}

// RsaDecrypt rsa decryption, parameter input untranscoded cipher text
func RsaDecrypt(privateKey []byte, cipherData []byte, opts ...RsaOption) ([]byte, error) {
	o := defaultRsaOptions()
	o.apply(opts...)

	return rsaDecryptWithPrivateKey(privateKey, cipherData, o.format)
}

// RsaEncryptHex rsa encryption, return hex
func RsaEncryptHex(publicKey []byte, rawData []byte, opts ...RsaOption) (string, error) {
	o := defaultRsaOptions()
	o.apply(opts...)

	cipherData, err := rsaEncryptWithPublicKey(publicKey, rawData)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(cipherData), nil
}

// RsaDecryptHex rsa decryption, return to original
func RsaDecryptHex(privateKey []byte, cipherHex string, opts ...RsaOption) (string, error) {
	o := defaultRsaOptions()
	o.apply(opts...)

	cipherData, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", err
	}

	rawData, err := rsaDecryptWithPrivateKey(privateKey, cipherData, o.format)
	if err != nil {
		return "", err
	}

	return string(rawData), nil
}

// RsaSign rsa signature, the returned ciphertext is not transcoded
func RsaSign(privateKey []byte, rawData []byte, opts ...RsaOption) ([]byte, error) {
	o := defaultRsaOptions()
	o.apply(opts...)

	return rsaSignWithPrivateKey(privateKey, o.hashType, rawData, o.format)
}

// RsaVerify rsa signature verification
func RsaVerify(publicKey []byte, rawData []byte, signData []byte, opts ...RsaOption) error {
	o := defaultRsaOptions()
	o.apply(opts...)

	return rsaVerifyWithPublicKey(publicKey, o.hashType, rawData, signData)
}

// RsaSignBase64 rsa signature, return base64
func RsaSignBase64(privateKey []byte, rawData []byte, opts ...RsaOption) (string, error) {
	o := defaultRsaOptions()
	o.apply(opts...)

	cipherData, err := rsaSignWithPrivateKey(privateKey, o.hashType, rawData, o.format)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(cipherData), nil
}

// RsaVerifyBase64 rsa signature verification
func RsaVerifyBase64(publicKey []byte, rawData []byte, signBase64 string, opts ...RsaOption) error {
	o := defaultRsaOptions()
	o.apply(opts...)

	signData, err := base64.StdEncoding.DecodeString(signBase64)
	if err != nil {
		return err
	}

	return rsaVerifyWithPublicKey(publicKey, o.hashType, rawData, signData)
}

// encrypt with public key
func rsaEncryptWithPublicKey(publicKey []byte, rawData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key is not pem format")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	prk, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("it's not a public key")
	}

	return rsa.EncryptPKCS1v15(rand.Reader, prk, rawData)
}

// decrypt with private key
func rsaDecryptWithPrivateKey(privateKey []byte, cipherData []byte, format string) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key is not pem format")
	}

	prk, err := getPrivateKey(block.Bytes, format)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, prk, cipherData)
}

// sign with private key
func rsaSignWithPrivateKey(privateKey []byte, hash crypto.Hash, rawData []byte, format string) ([]byte, error) {
	if !hash.Available() {
		return nil, errors.New("not supported hash type")
	}

	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key is not pem format")
	}

	prk, err := getPrivateKey(block.Bytes, format)
	if err != nil {
		return nil, err
	}

	h := hash.New()
	_, err = h.Write(rawData)
	if err != nil {
		return nil, err
	}
	hashed := h.Sum(nil)

	return rsa.SignPKCS1v15(rand.Reader, prk, hash, hashed)
}

// verify with public key
func rsaVerifyWithPublicKey(publicKey []byte, hash crypto.Hash, rawData []byte, signData []byte) (err error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return errors.New("public key is not pem format")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	prk, ok := pub.(*rsa.PublicKey)
	if !ok {
		return errors.New("it's not a public key")
	}

	h := hash.New()
	_, err = h.Write(rawData)
	if err != nil {
		return err
	}
	hashed := h.Sum(nil)

	return rsa.VerifyPKCS1v15(prk, hash, hashed, signData)
}

func getPrivateKey(der []byte, format string) (*rsa.PrivateKey, error) {
	var prk *rsa.PrivateKey
	switch format {
	case pkcs1:
		var err error
		prk, err = x509.ParsePKCS1PrivateKey(der)
		if err != nil {
			return nil, err
		}

	case pkcs8:
		priv, err := x509.ParsePKCS8PrivateKey(der)
		if err != nil {
			return nil, err
		}
		var ok bool
		prk, ok = priv.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("it's not a private key")
		}

	default:
		return nil, errors.New("unknown format = " + format)
	}

	return prk, nil
}
