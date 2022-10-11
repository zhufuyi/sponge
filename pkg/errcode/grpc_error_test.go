package errcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGRPCStatus(t *testing.T) {
	st := NewGRPCStatus(101, "something is wrong")
	err := st.Err()
	assert.Error(t, err)
	err = st.Err(Any("foo", "bar"))
	assert.Error(t, err)

	defer func() {
		recover()
	}()
	NewGRPCStatus(101, "something is wrong")
}

func TestToRPCCode(t *testing.T) {
	status := []*GRPCStatus{
		StatusSuccess,
		StatusInvalidParams,
		StatusUnauthorized,
		StatusInternalServerError,
		StatusNotFound,
		StatusAlreadyExists,
		StatusTimeout,
		StatusTooManyRequests,
		StatusForbidden,
		StatusLimitExceed,
		StatusDeadlineExceeded,
		StatusAccessDenied,
		StatusMethodNotAllowed,
	}

	var codes []string
	for _, s := range status {
		codes = append(codes, ToRPCCode(s.status.Code()).String())
	}
	t.Log(codes)
}

func TestGCode(t *testing.T) {
	code := GCode(1)
	t.Log("error code is", int(code))

	defer func() {
		recover()
	}()
	code = GCode(10001)
}
