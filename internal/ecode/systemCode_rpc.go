package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// rpc system level error code, with status prefix, error code range 30000~40000
var (
	StatusSuccess = errcode.StatusSuccess

	StatusCanceled            = errcode.StatusCanceled
	StatusUnknown             = errcode.StatusUnknown
	StatusInvalidParams       = errcode.StatusInvalidParams
	StatusDeadlineExceeded    = errcode.StatusDeadlineExceeded
	StatusNotFound            = errcode.StatusNotFound
	StatusAlreadyExists       = errcode.StatusAlreadyExists
	StatusPermissionDenied    = errcode.StatusPermissionDenied
	StatusResourceExhausted   = errcode.StatusResourceExhausted
	StatusFailedPrecondition  = errcode.StatusFailedPrecondition
	StatusAborted             = errcode.StatusAborted
	StatusOutOfRange          = errcode.StatusOutOfRange
	StatusUnimplemented       = errcode.StatusUnimplemented
	StatusInternalServerError = errcode.StatusInternalServerError
	StatusServiceUnavailable  = errcode.StatusServiceUnavailable
	StatusDataLoss            = errcode.StatusDataLoss
	StatusUnauthorized        = errcode.StatusUnauthorized

	StatusTimeout          = errcode.StatusTimeout
	StatusTooManyRequests  = errcode.StatusTooManyRequests
	StatusForbidden        = errcode.StatusForbidden
	StatusLimitExceed      = errcode.StatusLimitExceed
	StatusMethodNotAllowed = errcode.StatusMethodNotAllowed
	StatusAccessDenied     = errcode.StatusAccessDenied
)

// Any kev-value
func Any(key string, val interface{}) errcode.Detail {
	return errcode.Any(key, val)
}
