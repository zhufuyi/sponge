package errcode

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var errCodes = map[int]*Error{}

// Error 错误
type Error struct {
	// 错误码
	code int
	// 错误消息
	msg string
	// 详细信息
	details []string
}

// NewError 创建新错误信息
func NewError(code int, msg string) *Error {
	if v, ok := errCodes[code]; ok {
		panic(fmt.Sprintf("http error code = %d already exists, please replace with a new error code, old msg = %s", code, v.Msg()))
	}
	e := &Error{code: code, msg: msg}
	errCodes[code] = e
	return e
}

// Err 转为标准error
func (e *Error) Err() error {
	if len(e.details) == 0 {
		return fmt.Errorf("code = %d, msg = %s", e.code, e.msg)
	}
	return fmt.Errorf("code = %d, msg = %s, details = %v", e.code, e.msg, e.details)
}

// Code 错误码
func (e *Error) Code() int {
	return e.code
}

// Msg 错误信息
func (e *Error) Msg() string {
	return e.msg
}

// Msgf 附加信息
func (e *Error) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.msg, args...)
}

// Details 错误详情
func (e *Error) Details() []string {
	return e.details
}

// WithDetails 携带附加错误详情
func (e *Error) WithDetails(details ...string) *Error {
	newError := *e
	newError.details = []string{}
	newError.details = append(newError.details, details...)

	return &newError
}

// ToHTTPCode 转换为http错误码
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
	case Forbidden.Code():
		return http.StatusForbidden
	case NotFound.Code():
		return http.StatusNotFound
	case Timeout.Code():
		return http.StatusRequestTimeout
	}

	return e.Code()
}

// ParseError 根据标准错误信息解析出错误码和错误信息
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
