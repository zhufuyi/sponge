// nolint

package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// userExample rpc service level error code
// each resource name corresponds to a unique number (rpc type), the number range is 1~100, if there is the same number, trigger panic
var (
	_userExampleNO       = 1
	_userExampleName     = "userExample"
	_userExampleBaseCode = errcode.RCode(_userExampleNO)

	StatusCreateUserExample = errcode.NewRPCStatus(_userExampleBaseCode+1, "failed to create "+_userExampleName)
	StatusDeleteUserExample = errcode.NewRPCStatus(_userExampleBaseCode+2, "failed to delete "+_userExampleName)
	StatusUpdateUserExample = errcode.NewRPCStatus(_userExampleBaseCode+3, "failed to update "+_userExampleName)
	StatusGetUserExample    = errcode.NewRPCStatus(_userExampleBaseCode+4, "failed to get "+_userExampleName+" details")
	StatusListUserExample   = errcode.NewRPCStatus(_userExampleBaseCode+5, "failed to get list of "+_userExampleName)
	// for each error code added, add +1 to the previous error code
)
