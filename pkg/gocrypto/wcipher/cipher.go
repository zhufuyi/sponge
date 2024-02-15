// Package wcipher is a package to encrypt and decrypt data.
package wcipher

import (
	"crypto/cipher"
)

// Cipher provides a unified interface to encrypt/decrypt data.
type Cipher interface {
	Encrypt(src []byte) []byte
	Decrypt(src []byte) []byte
}

type blockCipher struct {
	padding Padding
	encrypt cipher.BlockMode
	decrypt cipher.BlockMode
}

// NewBlockCipher new block encryption
func NewBlockCipher(padding Padding, encrypt, decrypt cipher.BlockMode) Cipher {
	return &blockCipher{
		encrypt: encrypt,
		decrypt: decrypt,
		padding: padding,
	}
}

// Encrypt encrypted
func (blockCipher *blockCipher) Encrypt(plaintext []byte) []byte {
	plaintext = blockCipher.padding.Padding(plaintext, blockCipher.encrypt.BlockSize())
	ciphertext := make([]byte, len(plaintext))
	blockCipher.encrypt.CryptBlocks(ciphertext, plaintext)
	return ciphertext
}

// Decrypt decrypt
func (blockCipher *blockCipher) Decrypt(ciphertext []byte) []byte {
	plaintext := make([]byte, len(ciphertext))
	blockCipher.decrypt.CryptBlocks(plaintext, ciphertext)
	plaintext = blockCipher.padding.UnPadding(plaintext)
	return plaintext
}

// ---------------------------------------------------------------------------------------

type streamCipher struct {
	encrypt cipher.Stream
	decrypt cipher.Stream
}

// NewStreamCipher new stream encryption
func NewStreamCipher(encrypt cipher.Stream, decrypt cipher.Stream) Cipher {
	return &streamCipher{
		encrypt: encrypt,
		decrypt: decrypt}
}

// Encrypt stream encryption
func (streamCipher *streamCipher) Encrypt(plaintext []byte) []byte {
	ciphertext := make([]byte, len(plaintext))
	streamCipher.encrypt.XORKeyStream(ciphertext, plaintext)
	return ciphertext
}

// Decrypt stream decryption
func (streamCipher *streamCipher) Decrypt(ciphertext []byte) []byte {
	plaintext := make([]byte, len(ciphertext))
	streamCipher.decrypt.XORKeyStream(plaintext, ciphertext)
	return plaintext
}
