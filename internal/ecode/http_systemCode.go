// nolint

package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// http系统级别错误码，无Err前缀，错误码范围100000~200000
var (
	Success             = errcode.Success
	InvalidParams       = errcode.InvalidParams
	Unauthorized        = errcode.Unauthorized
	InternalServerError = errcode.InternalServerError
	NotFound            = errcode.NotFound
	AlreadyExists       = errcode.AlreadyExists
	Timeout             = errcode.Timeout
	TooManyRequests     = errcode.TooManyRequests
	Forbidden           = errcode.Forbidden
	LimitExceed         = errcode.LimitExceed

	DeadlineExceeded = errcode.DeadlineExceeded
	AccessDenied     = errcode.AccessDenied
	MethodNotAllowed = errcode.MethodNotAllowed
)
