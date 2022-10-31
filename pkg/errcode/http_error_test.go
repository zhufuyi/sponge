package errcode

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	code := 101
	msg := "something is wrong"

	e := NewError(code, msg)
	assert.Equal(t, code, e.Code())
	assert.Equal(t, msg, e.Msg())
	assert.Contains(t, e.Err().Error(), msg)
	assert.Contains(t, e.Msgf([]interface{}{"foo", "bar"}), msg)
	details := []string{"a", "b", "c"}
	assert.Equal(t, details, e.WithDetails(details...).Details())

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
	code = HCode(10001)
}
