// Package errcode is used for http and grpc error codes, include system-level error codes and business-level error codes
package errcode

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// ToHTTPCodeLabel need to convert to standard http code label
const ToHTTPCodeLabel = "[standard http code]"

var errCodes = map[int]*Error{}
var httpErrCodes = map[int]string{}

// Error error
type Error struct {
	code    int
	msg     string
	details []string

	// if true, need to convert to standard http code
	// use ErrToHTTP and ParseError will set this to true
	needHTTPCode bool
}

// NewError create a new error message
func NewError(code int, msg string, details ...string) *Error {
	if v, ok := errCodes[code]; ok {
		panic(fmt.Sprintf(`http error code = %d already exists, please define a new error code,
msg1 = %s
msg2 = %s
`, code, v.Msg(), msg))
	}

	httpErrCodes[code] = msg
	e := &Error{code: code, msg: msg, details: details}
	errCodes[code] = e
	return e
}

// Err convert to standard error,
// if there is a parameter 'msg', it will replace the original message.
func (e *Error) Err(msg ...string) error {
	message := e.msg
	if len(msg) > 0 {
		message = strings.Join(msg, ", ")
	}

	if len(e.details) == 0 {
		return fmt.Errorf("code = %d, msg = %s", e.code, message)
	}
	return fmt.Errorf("code = %d, msg = %s, details = %v", e.code, message, e.details)
}

// ErrToHTTP convert to standard error add ToHTTPCodeLabel to error message,
// use it if you need to convert to standard HTTP status code,
// if there is a parameter 'msg', it will replace the original message.
// Tips: you can call the GetErrorCode function to get the standard HTTP status code.
func (e *Error) ErrToHTTP(msg ...string) error {
	message := e.msg
	if len(msg) > 0 {
		message = strings.Join(msg, ", ")
	}

	if len(e.details) == 0 {
		return fmt.Errorf("code = %d, msg = %s%s", e.code, message, ToHTTPCodeLabel)
	}
	return fmt.Errorf("code = %d, msg = %s, details = %v%s", e.code, message, strings.Join(e.details, ", "), ToHTTPCodeLabel)
}

// Code get error code
func (e *Error) Code() int {
	return e.code
}

// Msg get error code message
func (e *Error) Msg() string {
	return e.msg
}

// NeedHTTPCode need to convert to standard http code
func (e *Error) NeedHTTPCode() bool {
	return e.needHTTPCode
}

// Details get error code details
func (e *Error) Details() []string {
	return e.details
}

// WithDetails add error details
func (e *Error) WithDetails(details ...string) *Error {
	newError := &Error{code: e.code, msg: e.msg}
	newError.msg += ", " + strings.Join(details, ", ")
	return newError
}

// RewriteMsg rewrite error message
func (e *Error) RewriteMsg(msg string) *Error {
	return &Error{code: e.code, msg: msg}
}

// WithOutMsg out error message
// Deprecated: use RewriteMsg instead
func (e *Error) WithOutMsg(msg string) *Error {
	return &Error{code: e.code, msg: msg}
}

// WithOutMsgI18n out error message i18n
// langMsg example:
//
//	map[int]map[string]string{
//			20010: {
//				"en-US": "login failed",
//				"zh-CN": "登录失败",
//			},
//		}
//
// lang BCP 47 code https://learn.microsoft.com/en-us/openspecs/office_standards/ms-oe376/6c085406-a698-4e12-9d4d-c3b0ee3dbc4a
func (e *Error) WithOutMsgI18n(langMsg map[int]map[string]string, lang string) *Error {
	if i18nMsg, ok := langMsg[e.Code()]; ok {
		if msg, ok2 := i18nMsg[lang]; ok2 {
			return &Error{code: e.code, msg: msg}
		}
	}

	return &Error{code: e.code, msg: e.msg}
}

// ToHTTPCode convert to http error code
func (e *Error) ToHTTPCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK
	case InternalServerError.Code():
		return http.StatusInternalServerError
	case InvalidParams.Code():
		return http.StatusBadRequest
	}

	switch e.Code() {
	case Unauthorized.Code(), PermissionDenied.Code():
		return http.StatusUnauthorized
	case TooManyRequests.Code(), LimitExceed.Code():
		return http.StatusTooManyRequests
	case Forbidden.Code(), AccessDenied.Code():
		return http.StatusForbidden
	case NotFound.Code():
		return http.StatusNotFound
	case Conflict.Code(), AlreadyExists.Code():
		return http.StatusConflict
	case TooEarly.Code():
		return http.StatusTooEarly
	case Timeout.Code(), DeadlineExceeded.Code():
		return http.StatusRequestTimeout
	case MethodNotAllowed.Code():
		return http.StatusMethodNotAllowed
	case ServiceUnavailable.Code():
		return http.StatusServiceUnavailable
	case Unimplemented.Code():
		return http.StatusNotImplemented
	case StatusBadGateway.Code():
		return http.StatusBadGateway
	}

	return http.StatusInternalServerError
}

// ParseError parsing out error codes from error messages
func ParseError(err error) *Error {
	if err == nil {
		return Success
	}

	outError := &Error{
		code: -1,
		msg:  "unknown error",
	}

	splits := strings.Split(err.Error(), ", msg = ")
	codeStr := strings.ReplaceAll(splits[0], "code = ", "")
	code, er := strconv.Atoi(codeStr)
	if er != nil {
		return outError
	}

	if e, ok := errCodes[code]; ok {
		if len(splits) > 1 {
			outError.code = code
			outError.msg = splits[1]
			outError.needHTTPCode = strings.Contains(err.Error(), ToHTTPCodeLabel)
			return outError
		}
		return e
	}

	return outError
}

// GetErrorCode get Error code from error returned by http invoke
func GetErrorCode(err error) int {
	e := ParseError(err)
	if e.needHTTPCode {
		return e.ToHTTPCode()
	}
	return e.Code()
}

// ListHTTPErrCodes list http error codes
func ListHTTPErrCodes() []ErrInfo {
	return getErrorInfo(httpErrCodes)
}
