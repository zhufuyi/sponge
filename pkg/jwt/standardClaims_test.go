package jwt

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateTokenStandard(t *testing.T) {
	opt = nil
	token, err := GenerateTokenStandard()
	assert.Error(t, err)

	Init()
	token, err = GenerateTokenStandard()
	assert.NoError(t, err)
	t.Log(token)
}

func TestVerifyTokenStandard(t *testing.T) {
	opt = nil
	err := VerifyTokenStandard("token")
	assert.Error(t, err)

	Init(WithSigningKey("123456"))

	// normal verify
	token, err := GenerateTokenStandard()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(token)
	err = VerifyTokenStandard(token)
	if err != nil {
		t.Fatal(err)
	}

	// invalid token format
	token2 := "xxx.xxx.xxx"
	err = VerifyTokenStandard(token2)
	assert.Equal(t, err, errFormat)

	// signature failure
	token3 := token + "xxx"
	err = VerifyTokenStandard(token3)
	assert.Equal(t, err, errSignature)

	// token has expired
	Init(
		WithSigningKey("123456"),
		WithExpire(time.Second),
	)
	token, err = GenerateTokenStandard()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 2)
	err = VerifyTokenStandard(token)
	assert.Equal(t, err, errExpired)
}
