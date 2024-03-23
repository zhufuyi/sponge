package jwt

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	opt = nil
	token, err := GenerateToken("123")
	assert.Error(t, err)

	Init()
	token, err = GenerateToken("123")
	assert.NoError(t, err)

	v, err := ParseToken(token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)

	time.Sleep(time.Second)
	newToken, err := RefreshToken(token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token, newToken)
}

func TestParseToken(t *testing.T) {
	opt = nil
	v, err := ParseToken("token")
	assert.Error(t, err)

	uid := "123"
	name := "admin"

	Init(
		WithSigningKey("123456"),
		WithExpire(time.Second),
		WithSigningMethod(HS512),
	)

	// success
	token, err := GenerateToken(uid, name)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(token)
	v, err = ParseToken(token)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(v)

	// invalid token format
	token2 := "xxx.xxx.xxx"
	v, err = ParseToken(token2)
	assert.Error(t, err)

	// signature failure
	token3 := token + "xxx"
	v, err = ParseToken(token3)
	assert.Error(t, err)

	// token has expired
	token, err = GenerateToken(uid, name)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)
	v, err = ParseToken(token)
	assert.True(t, errors.Is(err, ErrTokenExpired))
}

func TestGenerateCustomToken(t *testing.T) {
	fields := KV{"id": 123, "foo": "bar"}
	Init()
	token, err := GenerateCustomToken(fields)
	assert.NoError(t, err)

	claims, err := ParseCustomToken(token)
	assert.NoError(t, err)
	uid, _ := claims.Get("id")
	assert.NotNil(t, uid)
	foo, _ := claims.Get("foo")
	assert.NotNil(t, foo)
	t.Log(uid, foo)

	claims.Fields = nil
	foo, _ = claims.Get("foo")
	assert.Nil(t, foo)

	time.Sleep(time.Second)
	newToken, err := RefreshCustomToken(token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token, newToken)
}

func TestParseCustomToken(t *testing.T) {
	fields := KV{"id": 123, "foo": "bar"}
	opt = nil
	v, err := ParseCustomToken("token")
	assert.Error(t, err)

	Init(
		WithSigningKey("123456"),
		WithExpire(time.Second),
		WithSigningMethod(HS512),
	)

	// success
	token, err := GenerateCustomToken(fields)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(token)
	v, err = ParseCustomToken(token)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(v)

	// invalid token format
	token2 := "xxx.xxx.xxx"
	v, err = ParseCustomToken(token2)
	assert.Error(t, err)

	// signature failure
	token3 := token + "xxx"
	v, err = ParseCustomToken(token3)
	assert.Error(t, err)

	// token has expired
	token, err = GenerateCustomToken(fields)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)
	v, err = ParseCustomToken(token)
	assert.True(t, errors.Is(err, ErrTokenExpired))
}
