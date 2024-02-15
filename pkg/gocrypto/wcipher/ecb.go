package wcipher

import (
	"crypto/cipher"
)

type ecb struct {
	block     cipher.Block
	blockSize int
}

type ecbEncrypt ecb

func (e *ecbEncrypt) BlockSize() int {
	return e.blockSize
}

func (e *ecbEncrypt) CryptBlocks(dst, src []byte) {
	if len(src)%e.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	for len(src) > 0 {
		e.block.Encrypt(dst, src[:e.blockSize])
		src = src[e.blockSize:]
		dst = dst[e.blockSize:]
	}
}

type ecbDecrypt ecb

func (e *ecbDecrypt) BlockSize() int {
	return e.blockSize
}

func (e *ecbDecrypt) CryptBlocks(dst, src []byte) {
	if len(src)%e.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	for len(src) > 0 {
		e.block.Decrypt(dst, src[:e.blockSize])
		src = src[e.blockSize:]
		dst = dst[e.blockSize:]
	}
}

// NewECBEncrypt ecb encrypt
func NewECBEncrypt(block cipher.Block) cipher.BlockMode {
	return &ecbEncrypt{block: block, blockSize: block.BlockSize()}
}

// NewECBDecrypt ecb decrypt
func NewECBDecrypt(block cipher.Block) cipher.BlockMode {
	return &ecbDecrypt{block: block, blockSize: block.BlockSize()}
}
