package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// nolint
// rpc系统级别错误码，有status前缀
var (
	StatusSuccess = errcode.StatusSuccess

	StatusInvalidParams       = errcode.StatusInvalidParams
	StatusUnauthorized        = errcode.StatusUnauthorized
	StatusInternalServerError = errcode.StatusInternalServerError
	StatusNotFound            = errcode.StatusNotFound
	StatusAlreadyExists       = errcode.StatusAlreadyExists
	StatusTimeout             = errcode.StatusTimeout
	StatusTooManyRequests     = errcode.StatusTooManyRequests
	StatusForbidden           = errcode.StatusForbidden
	StatusLimitExceed         = errcode.StatusLimitExceed

	StatusDeadlineExceeded = errcode.StatusDeadlineExceeded
	StatusAccessDenied     = errcode.StatusAccessDenied
	StatusMethodNotAllowed = errcode.StatusMethodNotAllowed
)
