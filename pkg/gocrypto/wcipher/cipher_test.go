package wcipher

import (
	"crypto/aes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var key = []byte("0123456789012345")
var iv = key
var text = []byte("123456")

func TestBlockCipher(t *testing.T) {
	p := NewPKCS57Padding()
	b, _ := aes.NewCipher(key)

	cm := NewECBMode()
	cm.SetPadding(p)
	c := cm.Cipher(b, iv)
	eText := c.Encrypt(text)
	dText := c.Decrypt(eText)
	assert.Equal(t, text, dText)

	cm = NewCFBMode()
	cm.SetPadding(p)
	c = cm.Cipher(b, iv)
	eText = c.Encrypt(text)
	dText = c.Decrypt(eText)
	assert.Equal(t, text, dText)

	cm = NewCBCMode()
	cm.SetPadding(p)
	c = cm.Cipher(b, iv)
	eText = c.Encrypt(text)
	dText = c.Decrypt(eText)
	assert.Equal(t, text, dText)

	e := &ecbDecrypt{
		blockSize: 16,
	}
	t.Log(e.BlockSize())
}

func TestNewAES(t *testing.T) {
	c, _ := NewAES(key)
	eText := c.Encrypt(text)
	dText := c.Decrypt(eText)
	assert.Equal(t, text, dText)

	c, _ = NewAESWith(key, NewECBMode())
	eText = c.Encrypt(text)
	dText = c.Decrypt(eText)
	assert.Equal(t, text, dText)

	c, _ = NewDES(key[:8])
	eText = c.Encrypt(text)
	dText = c.Decrypt(eText)
	assert.Equal(t, text, dText)

	c, _ = NewDESWith(key[:8], NewCBCMode())
	eText = c.Encrypt(text)
	dText = c.Decrypt(eText)
	assert.Equal(t, text, dText)

	c, _ = NewDESWith(key[:8], NewCTRMode())
	eText = c.Encrypt(text)
	dText = c.Decrypt(eText)
	assert.Equal(t, text, dText)

	// test err
	_, err := NewAES(key[:8])
	assert.NotNil(t, err)
	_, err = NewAESWith(key[:8], NewCFBMode())
	assert.NotNil(t, err)
	_, err = NewDES(key)
	assert.NotNil(t, err)
	_, err = NewDESWith(key, NewCFBMode())
	assert.NotNil(t, err)
}
