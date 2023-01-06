package errcode

import (
	"errors"
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
	assert.Contains(t, e.Err().Error(), msg)
	details := []string{"a", "b", "c"}
	assert.Contains(t, e.WithDetails(details...).Err().Error(), strings.Join(details, ", "))

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
		NewError(1010, "unknown"),
	}

	var httpCodes []int
	for _, ec := range errorsCodes {
		httpCodes = append(httpCodes, ec.ToHTTPCode())
	}
	t.Log(httpCodes)

	var codes []int
	for _, ec := range errorsCodes {
		e := ParseError(ec.Err())
		codes = append(codes, e.Code())
	}
	e = ParseError(errors.New("unknown error"))
	codes = append(codes, e.Code())
	t.Log(codes)

	_ = ParseError(nil)
	_ = e.Details()

	defer func() {
		recover()
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
