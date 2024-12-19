package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// userExample business-level http error codes.
// the userExampleNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	userExampleNO       = 1
	userExampleName     = "userExample"
	userExampleBaseCode = errcode.HCode(userExampleNO)

	ErrCreateUserExample     = errcode.NewError(userExampleBaseCode+1, "failed to create "+userExampleName)
	ErrDeleteByIDUserExample = errcode.NewError(userExampleBaseCode+2, "failed to delete "+userExampleName)
	ErrUpdateByIDUserExample = errcode.NewError(userExampleBaseCode+3, "failed to update "+userExampleName)
	ErrGetByIDUserExample    = errcode.NewError(userExampleBaseCode+4, "failed to get "+userExampleName+" details")
	ErrListUserExample       = errcode.NewError(userExampleBaseCode+5, "failed to list of "+userExampleName)

	// error codes are globally unique, adding 1 to the previous error code
)
