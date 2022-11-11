// nolint

package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// userExample rpc服务级别错误码
// 每个资源名称对应唯一编号(rpc类型)，编号范围1~100，如果存在编号相同，触发panic
var (
	_userExampleNO       = 1
	_userExampleName     = "userExample"
	_userExampleBaseCode = errcode.RCode(_userExampleNO)

	StatusCreateUserExample = errcode.NewRPCStatus(_userExampleBaseCode+1, "failed to create "+_userExampleName)
	StatusDeleteUserExample = errcode.NewRPCStatus(_userExampleBaseCode+2, "failed to delete "+_userExampleName)
	StatusUpdateUserExample = errcode.NewRPCStatus(_userExampleBaseCode+3, "failed to update "+_userExampleName)
	StatusGetUserExample    = errcode.NewRPCStatus(_userExampleBaseCode+4, "failed to get "+_userExampleName+" details")
	StatusListUserExample   = errcode.NewRPCStatus(_userExampleBaseCode+5, "failed to get list of "+_userExampleName)
	// 每添加一个错误码，在上一个错误码基础上+1
)
