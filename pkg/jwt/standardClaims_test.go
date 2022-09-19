package jwt

import (
	"fmt"
	"testing"
)

func TestVerifyTokenStandard(t *testing.T) {
	Init(WithSigningKey("123456"))

	token, err := GenerateTokenStandard()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(token)

	err = VerifyTokenStandard(token)
	if err != nil {
		t.Error(err)
	}
}
