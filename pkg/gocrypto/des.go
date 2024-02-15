// Symmetric encryption DES, one of the most popular encryption algorithms,
// is gradually being replaced by AES.

package gocrypto

import (
	"encoding/hex"

	"testCode/pkg/gocrypto/wcipher"
)

// DesEncrypt des encryption, the returned ciphertext is not transcoded
func DesEncrypt(rawData []byte, opts ...DesOption) ([]byte, error) {
	o := defaultDesOptions()
	o.apply(opts...)

	return desEncrypt(o.mode, rawData, o.desKey)
}

// DesDecrypt des decryption, parameter input untranscoded cipher text
func DesDecrypt(cipherData []byte, opts ...DesOption) ([]byte, error) {
	o := defaultDesOptions()
	o.apply(opts...)

	return desDecrypt(o.mode, cipherData, o.desKey)
}

// DesEncryptHex des encrypts and returns a ciphertext that has been transcoded
func DesEncryptHex(rawData string, opts ...DesOption) (string, error) {
	o := defaultDesOptions()
	o.apply(opts...)

	cipherData, err := desEncrypt(o.mode, []byte(rawData), o.desKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(cipherData), nil
}

// DesDecryptHex des decryption, parameter input has been transcoded ciphertext string
func DesDecryptHex(cipherStr string, opts ...DesOption) (string, error) {
	o := defaultDesOptions()
	o.apply(opts...)

	cipherData, err := hex.DecodeString(cipherStr)
	if err != nil {
		return "", err
	}

	rawData, err := desDecrypt(o.mode, cipherData, o.desKey)
	if err != nil {
		return "", err
	}

	return string(rawData), nil
}

func desEncrypt(mode string, rawData []byte, key []byte) ([]byte, error) {
	cipherMode, err := getCipherMode(mode)
	if err != nil {
		return nil, err
	}

	cip, err := wcipher.NewDESWith(key, cipherMode)
	if err != nil {
		return nil, err
	}

	return cip.Encrypt(rawData), nil
}

func desDecrypt(mode string, cipherData []byte, key []byte) ([]byte, error) {
	cipherMode, err := getCipherMode(mode)
	if err != nil {
		return nil, err
	}

	cip, err := wcipher.NewDESWith(key, cipherMode)
	if err != nil {
		return nil, err
	}

	return cip.Decrypt(cipherData), nil
}
