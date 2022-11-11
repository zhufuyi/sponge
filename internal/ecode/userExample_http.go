// nolint

package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// userExample http服务级别错误码
// 每个资源名称对应唯一编号(http类型)，编号范围1~100，如果存在编号相同，触发panic
var (
	userExampleNO       = 1
	userExampleName     = "userExample"
	userExampleBaseCode = errcode.HCode(userExampleNO)

	ErrCreateUserExample = errcode.NewError(userExampleBaseCode+1, "failed to create "+userExampleName)
	ErrDeleteUserExample = errcode.NewError(userExampleBaseCode+2, "failed to delete "+userExampleName)
	ErrUpdateUserExample = errcode.NewError(userExampleBaseCode+3, "failed to update "+userExampleName)
	ErrGetUserExample    = errcode.NewError(userExampleBaseCode+4, "failed to get "+userExampleName+" details")
	ErrListUserExample   = errcode.NewError(userExampleBaseCode+5, "failed to get list of "+userExampleName)
	// 每添加一个错误码，在上一个错误码基础上+1
)
