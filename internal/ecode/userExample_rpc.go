package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// userExample rpc service level error code
// each resource name corresponds to a unique number (rpc type), the number range is 1~100, if there is the same number, trigger panic
var (
	_userExampleNO       = 2
	_userExampleName     = "userExample"
	_userExampleBaseCode = errcode.RCode(_userExampleNO)

	StatusCreateUserExample      = errcode.NewRPCStatus(_userExampleBaseCode+1, "failed to create "+_userExampleName)
	StatusDeleteUserExample      = errcode.NewRPCStatus(_userExampleBaseCode+2, "failed to delete "+_userExampleName)
	StatusDeleteByIDsUserExample = errcode.NewRPCStatus(_userExampleBaseCode+3, "failed to delete by batch ids "+_userExampleName)
	StatusUpdateUserExample      = errcode.NewRPCStatus(_userExampleBaseCode+4, "failed to update "+_userExampleName)
	StatusGetUserExample         = errcode.NewRPCStatus(_userExampleBaseCode+5, "failed to get "+_userExampleName+" details")
	StatusListByIDsUserExample   = errcode.NewRPCStatus(_userExampleBaseCode+6, "failed to list by batch ids "+_userExampleName)
	StatusListUserExample        = errcode.NewRPCStatus(_userExampleBaseCode+7, "failed to list of "+_userExampleName)
	// error codes are globally unique, adding 1 to the previous error code
)
