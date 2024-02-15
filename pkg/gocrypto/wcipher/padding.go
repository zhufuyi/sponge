package wcipher

import (
	"bytes"
)

// Padding provides a unified interface to populate/restore data for various population methods.
type Padding interface {
	Padding(src []byte, blockSize int) []byte
	UnPadding(src []byte) []byte
}

type padding struct{}

type pkcs57Padding padding

// NewPKCS57Padding new pkcs57 padding
func NewPKCS57Padding() Padding {
	return &pkcs57Padding{}
}

// Padding pkcs57 padding
func (p *pkcs57Padding) Padding(src []byte, blockSize int) []byte {
	paddingSize := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)
	return append(src, padText...)
}

// UnPadding pkcs57 un-padding
func (p *pkcs57Padding) UnPadding(src []byte) []byte {
	length := len(src)
	unPadding := int(src[length-1])
	return src[:(length - unPadding)]
}
