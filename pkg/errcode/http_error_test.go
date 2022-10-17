package errcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	code := 101
	msg := "something is wrong"

	e := NewError(code, msg)
	assert.Equal(t, code, e.Code())
	assert.Equal(t, msg, e.Msg())
	assert.Contains(t, e.Error(), msg)
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
	for _, e := range errorsCodes {
		httpCodes = append(httpCodes, e.ToHTTPCode())
	}
	t.Log(httpCodes)

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
