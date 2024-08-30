package errcode

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errorsCodes = []*Error{
	Success,
	InvalidParams,
	Unauthorized,
	InternalServerError,
	NotFound,
	Conflict,
	AlreadyExists,
	Timeout,
	TooManyRequests,
	Forbidden,
	LimitExceed,
	DeadlineExceeded,
	AccessDenied,
	MethodNotAllowed,
	ServiceUnavailable,
	TooEarly,

	Canceled,
	Unknown,
	PermissionDenied,
	ResourceExhausted,
	FailedPrecondition,
	Aborted,
	OutOfRange,
	Unimplemented,
	StatusBadGateway,
}

func TestNewError(t *testing.T) {
	code := 21101
	msg := "something is wrong"

	e := NewError(code, msg)
	assert.Equal(t, code, e.Code())
	assert.Equal(t, msg, e.Msg())
	assert.Equal(t, false, e.NeedHTTPCode())
	assert.Contains(t, e.Err().Error(), msg)
	assert.Contains(t, e.ErrToHTTP().Error(), ToHTTPCodeLabel)
	details := []string{"a", "b", "c"}
	assert.Contains(t, e.WithDetails(details...).Err().Error(), strings.Join(details, ", "))
	assert.Contains(t, e.WithDetails(details...).ErrToHTTP().Error(), ToHTTPCodeLabel)

	errorsCodes = append(errorsCodes,
		DataLoss.WithDetails("foo", "bar"),
		DataLoss.WithOutMsg("foobar"),
		NewError(1010, "unknown"),
	)

	var httpCodes []int
	for _, ec := range errorsCodes {
		httpCodes = append(httpCodes, ec.ToHTTPCode())
	}
	t.Log(httpCodes)

	defer func() {
		if err := recover(); err != nil {
			t.Log(err)
		}
	}()
	_ = NewError(code, msg)
}

func TestHCode(t *testing.T) {
	code := HCode(1)
	t.Log("error code is", code)

	defer func() {
		recover()
	}()
	code = HCode(101)
}

func TestListHTTPErrCodes(t *testing.T) {
	errInfos := ListHTTPErrCodes()
	for _, v := range errInfos {
		fmt.Println(v.Code, v.Msg)
	}
}

func TestParseError(t *testing.T) {
	errorsCodes = append(errorsCodes,
		NewError(21102, "something is wrong"),
		ParseError(errors.New("unknown error")),
	)

	var relationshipCodes []string
	for _, ec := range errorsCodes {
		e1 := ParseError(ec.Err())
		e2 := ParseError(ec.ErrToHTTP())
		relationshipCodes = append(relationshipCodes, fmt.Sprintf("%d:%d", e1.Code(), e2.ToHTTPCode()))
	}
	t.Log(relationshipCodes)

	e := ParseError(nil)
	t.Log(e)
}

func TestGetErrorCode(t *testing.T) {
	for _, e := range errorsCodes {
		t.Log(e.Code(), "|",
			GetErrorCode(e.Err()),
			GetErrorCode(e.Err("reason for error")), "|",

			GetErrorCode(e.ErrToHTTP()),
			GetErrorCode(e.ErrToHTTP("reason for error")),
		)
	}
}

func TestError_WithOutMsgI18n(t *testing.T) {
	var langMsg = map[int]map[string]string{
		20011: {
			"en-US": "login failed",
			"zh-CN": "登录失败",
		},
	}

	e := NewError(20011, "login failed")
	e1 := e.WithOutMsgI18n(langMsg, "zh-CN")
	assert.Equal(t, "登录失败", e1.Msg())

	e2 := e.WithOutMsgI18n(langMsg, "zh")
	assert.NotEqual(t, "登录失败", e2.Msg())

	t.Log(e1.Msg(), e2.Msg())
}
