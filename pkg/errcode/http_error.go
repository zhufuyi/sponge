package errcode

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var errCodes = map[int]*Error{}

// Error error
type Error struct {
	code    int
	msg     string
	details []string
}

// NewError create a new error message
func NewError(code int, msg string) *Error {
	if v, ok := errCodes[code]; ok {
		panic(fmt.Sprintf("http error code = %d already exists, please replace with a new error code, old msg = %s", code, v.Msg()))
	}
	e := &Error{code: code, msg: msg}
	errCodes[code] = e
	return e
}

// Err convert to standard error
func (e *Error) Err() error {
	if len(e.details) == 0 {
		return fmt.Errorf("code = %d, msg = %s", e.code, e.msg)
	}
	return fmt.Errorf("code = %d, msg = %s, details = %v", e.code, e.msg, e.details)
}

// Code get error code
func (e *Error) Code() int {
	return e.code
}

// Msg get error code message
func (e *Error) Msg() string {
	return e.msg
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

// ToHTTPCode convert to http error code
func (e *Error) ToHTTPCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK
	case InternalServerError.Code():
		return http.StatusInternalServerError
	case InvalidParams.Code():
		return http.StatusBadRequest
	case Unauthorized.Code():
		return http.StatusUnauthorized
	case TooManyRequests.Code(), LimitExceed.Code():
		return http.StatusTooManyRequests
	case Forbidden.Code(), AccessDenied.Code():
		return http.StatusForbidden
	case NotFound.Code():
		return http.StatusNotFound
	case AlreadyExists.Code():
		return http.StatusConflict
	case Timeout.Code(), DeadlineExceeded.Code():
		return http.StatusRequestTimeout
	case MethodNotAllowed.Code():
		return http.StatusMethodNotAllowed
	case ServiceUnavailable.Code():
		return http.StatusServiceUnavailable
	}

	return e.Code()
}

// ParseError parsing out error codes from error messages
func ParseError(err error) *Error {
	if err == nil {
		return Success
	}

	unknownError := &Error{
		code: -1,
		msg:  "unknown error",
	}

	splits := strings.Split(err.Error(), ", msg = ")
	codeStr := strings.ReplaceAll(splits[0], "code = ", "")
	code, er := strconv.Atoi(codeStr)
	if er != nil {
		return unknownError
	}

	if e, ok := errCodes[code]; ok {
		return e
	}

	return unknownError
}
