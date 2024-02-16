// Symmetric encryption AES, the advanced encryption standard with
// the highest level of security, has gradually replaced DES as the new
// generation of symmetric encryption standard.

package gocrypto

import (
	"encoding/hex"
	"errors"

	"github.com/zhufuyi/sponge/pkg/gocrypto/wcipher"
)

// AesEncrypt aes encryption, returns ciphertext is not transcoded
func AesEncrypt(rawData []byte, opts ...AesOption) ([]byte, error) {
	o := defaultAesOptions()
	o.apply(opts...)

	return aesEncryptByMode(o.mode, rawData, o.aesKey)
}

// AesDecrypt aes decryption, parameter input un-transcode cipher text
func AesDecrypt(cipherData []byte, opts ...AesOption) ([]byte, error) {
	o := defaultAesOptions()
	o.apply(opts...)

	return aesDecryptByMode(o.mode, cipherData, o.aesKey)
}

// AesEncryptHex aes encryption, the returned ciphertext is transcoded
func AesEncryptHex(rawData string, opts ...AesOption) (string, error) {
	o := defaultAesOptions()
	o.apply(opts...)

	cipherData, err := aesEncryptByMode(o.mode, []byte(rawData), o.aesKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(cipherData), nil
}

// AesDecryptHex aes decryption, parameter input has been transcoded ciphertext string
func AesDecryptHex(cipherStr string, opts ...AesOption) (string, error) {
	o := defaultAesOptions()
	o.apply(opts...)

	cipherData, err := hex.DecodeString(cipherStr)
	if err != nil {
		return "", err
	}

	rawData, err := aesDecryptByMode(o.mode, cipherData, o.aesKey)
	if err != nil {
		return "", err
	}

	return string(rawData), nil
}

func getCipherMode(mode string) (wcipher.CipherMode, error) {
	var cipherMode wcipher.CipherMode
	switch mode {
	case modeECB:
		cipherMode = wcipher.NewECBMode()
	case modeCBC:
		cipherMode = wcipher.NewCBCMode()
	case modeCFB:
		cipherMode = wcipher.NewCFBMode()
	case modeCTR:
		cipherMode = wcipher.NewCTRMode()
	default:
		return nil, errors.New("unknown mode = " + mode)
	}

	return cipherMode, nil
}

func aesEncryptByMode(mode string, rawData []byte, key []byte) ([]byte, error) {
	cipherMode, err := getCipherMode(mode)
	if err != nil {
		return nil, err
	}

	cip, err := wcipher.NewAESWith(key, cipherMode)
	if err != nil {
		return nil, err
	}

	return cip.Encrypt(rawData), nil
}

func aesDecryptByMode(mode string, cipherData []byte, key []byte) ([]byte, error) {
	cipherMode, err := getCipherMode(mode)
	if err != nil {
		return nil, err
	}

	cip, err := wcipher.NewAESWith(key, cipherMode)
	if err != nil {
		return nil, err
	}

	return cip.Decrypt(cipherData), nil
}
