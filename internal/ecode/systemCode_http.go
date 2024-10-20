// Package ecode is the package that unifies the definition of http error codes or grpc error codes here.
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
	Timeout             = errcode.Timeout
	TooManyRequests     = errcode.TooManyRequests
	Forbidden           = errcode.Forbidden
	LimitExceed         = errcode.LimitExceed
	Conflict            = errcode.Conflict
	TooEarly            = errcode.TooEarly

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

var SkipResponse = errcode.SkipResponse

// GetErrorCode get error code from error
var GetErrorCode = errcode.GetErrorCode
