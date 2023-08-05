package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// userExample http service level error code
// each resource name corresponds to a unique number (http type), the number range is 1~100, if there is the same number, trigger panic
var (
	userExampleNO       = 1
	userExampleName     = "userExample"
	userExampleBaseCode = errcode.HCode(userExampleNO)

	ErrCreateUserExample      = errcode.NewError(userExampleBaseCode+1, "failed to create "+userExampleName)
	ErrDeleteUserExample      = errcode.NewError(userExampleBaseCode+2, "failed to delete "+userExampleName)
	ErrDeleteByIDsUserExample = errcode.NewError(userExampleBaseCode+3, "failed to delete by batch ids "+userExampleName)
	ErrUpdateUserExample      = errcode.NewError(userExampleBaseCode+4, "failed to update "+userExampleName)
	ErrGetUserExample         = errcode.NewError(userExampleBaseCode+5, "failed to get "+userExampleName+" details")
	ErrListByIDsUserExample   = errcode.NewError(userExampleBaseCode+6, "failed to list by batch ids "+userExampleName)
	ErrListUserExample        = errcode.NewError(userExampleBaseCode+7, "failed to list of "+userExampleName)
	// for each error code added, add +1 to the previous error code
)
