package wcipher

import (
	"crypto/aes"
	"crypto/des"
)

// NewAES Create default AES cipher, use ECB working mode, pkcs57 padding,
// algorithm secret key length 128 192 256 bits, use secret key as initial vector.
func NewAES(key []byte) (Cipher, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return NewECBMode().Cipher(block, key[:block.BlockSize()]), err
}

// NewAESWith According to the specified working mode, create AES cipher,
// the length of the algorithm secret key is 128 192 256 bits, and the secret
// key is used as the initial vector.
func NewAESWith(key []byte, mode CipherMode) (Cipher, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return mode.Cipher(block, key[:block.BlockSize()]), nil
}

// NewDES Create default DES cipher, use ECB working mode, pkcs57 padding,
// algorithm secret key length 64 bits, use secret key as initial vector.
func NewDES(key []byte) (Cipher, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return NewECBMode().Cipher(block, key[:block.BlockSize()]), nil
}

// NewDESWith According to the specified working mode, create DES cipher,
// the length of the algorithm secret key is 64 bits, and use the secret key as the initial vector.
func NewDESWith(key []byte, mode CipherMode) (Cipher, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return mode.Cipher(block, key[:block.BlockSize()]), nil
}
