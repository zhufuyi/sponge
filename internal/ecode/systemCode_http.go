// nolint

package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// http system level error code, error code range 10000~20000
var (
	Success = errcode.Success

	InvalidParams       = errcode.InvalidParams
	Unauthorized        = errcode.Unauthorized
	InternalServerError = errcode.InternalServerError
	NotFound            = errcode.NotFound
	AlreadyExists       = errcode.AlreadyExists
	Timeout             = errcode.Timeout
	TooManyRequests     = errcode.TooManyRequests
	Forbidden           = errcode.Forbidden
	LimitExceed         = errcode.LimitExceed

	DeadlineExceeded   = errcode.DeadlineExceeded
	AccessDenied       = errcode.AccessDenied
	MethodNotAllowed   = errcode.MethodNotAllowed
	ServiceUnavailable = errcode.ServiceUnavailable

	Canceled           = errcode.Canceled
	Unknown            = errcode.Unknown
	PermissionDenied   = errcode.PermissionDenied
	ResourceExhausted  = errcode.ResourceExhausted
	FailedPrecondition = errcode.FailedPrecondition
	Aborted            = errcode.Aborted
	OutOfRange         = errcode.OutOfRange
	Unimplemented      = errcode.Unimplemented
	DataLoss           = errcode.DataLoss
)
