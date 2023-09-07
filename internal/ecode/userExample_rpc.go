package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// userExample business-level rpc error codes.
// the _userExampleNO value range is 1~100, if the same number appears, it will cause a failure to start the service.
var (
	_userExampleNO       = 2
	_userExampleName     = "userExample"
	_userExampleBaseCode = errcode.RCode(_userExampleNO)

	StatusCreateUserExample      = errcode.NewRPCStatus(_userExampleBaseCode+1, "failed to create "+_userExampleName)
	StatusDeleteByIDUserExample  = errcode.NewRPCStatus(_userExampleBaseCode+2, "failed to delete "+_userExampleName)
	StatusDeleteByIDsUserExample = errcode.NewRPCStatus(_userExampleBaseCode+3, "failed to delete by batch ids "+_userExampleName)
	StatusUpdateByIDUserExample  = errcode.NewRPCStatus(_userExampleBaseCode+4, "failed to update "+_userExampleName)
	StatusGetByIDUserExample     = errcode.NewRPCStatus(_userExampleBaseCode+5, "failed to get "+_userExampleName+" details")
	StatusListByIDsUserExample   = errcode.NewRPCStatus(_userExampleBaseCode+6, "failed to list by batch ids "+_userExampleName)
	StatusListUserExample        = errcode.NewRPCStatus(_userExampleBaseCode+7, "failed to list of "+_userExampleName)
	// error codes are globally unique, adding 1 to the previous error code
)
