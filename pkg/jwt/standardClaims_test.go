package jwt

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateTokenStandard(t *testing.T) {
	Init()
	token, err := GenerateTokenStandard()
	assert.NoError(t, err)
	t.Log(token)
}

func TestVerifyTokenStandard(t *testing.T) {
	Init(WithSigningKey("123456"))

	// 正常验证
	token, err := GenerateTokenStandard()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(token)
	err = VerifyTokenStandard(token)
	if err != nil {
		t.Fatal(err)
	}

	// 无效token格式
	token2 := "xxx.xxx.xxx"
	err = VerifyTokenStandard(token2)
	if !compareErr(err, errFormat) {
		t.Fatal(err)
	}

	// 签名失败
	token3 := token + "xxx"
	err = VerifyTokenStandard(token3)
	if !compareErr(err, errSignature) {
		t.Fatal(err)
	}

	// token已过期
	Init(
		WithSigningKey("123456"),
		WithExpire(time.Millisecond*200),
	)
	token, err = GenerateTokenStandard()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	err = VerifyTokenStandard(token)
	if !compareErr(err, errExpired) {
		t.Fatal(err)
	}
}
