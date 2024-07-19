package errcode

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	errorsCodes := []*Error{
		Success,
		InvalidParams,
		Unauthorized,
		InternalServerError,
		NotFound,
		AlreadyExists,
		Timeout,
		TooManyRequests,
		Forbidden,
		LimitExceed,
		DeadlineExceeded,
		AccessDenied,
		MethodNotAllowed,
		ServiceUnavailable,

		Canceled,
		Unknown,
		PermissionDenied,
		ResourceExhausted,
		FailedPrecondition,
		Aborted,
		OutOfRange,
		Unimplemented,
		StatusBadGateway,
		DataLoss.WithDetails("foo", "bar"),
		DataLoss.WithOutMsg("foobar"),
		NewError(1010, "unknown"),
	}

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
	errorsCodes := []*Error{
		Success,
		InvalidParams,
		Unauthorized,
		InternalServerError,
		NotFound,
		AlreadyExists,
		NewError(21102, "something is wrong"),
		ParseError(errors.New("unknown error")),
	}

	var codes1 []int
	var codes2 []int
	for _, ec := range errorsCodes {
		e1 := ParseError(ec.Err())
		codes1 = append(codes1, e1.Code())
		e2 := ParseError(ec.ErrToHTTP())
		codes2 = append(codes2, e2.ToHTTPCode())
	}
	t.Log(codes1, codes2)

	e := ParseError(nil)
	t.Log(e)
}
