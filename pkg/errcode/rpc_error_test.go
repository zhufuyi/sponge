package errcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRPCStatus(t *testing.T) {
	st := NewRPCStatus(41101, "something is wrong")
	err := st.Err()
	assert.Error(t, err)
	err = st.Err(Any("foo", "bar"))
	assert.Error(t, err)

	defer func() {
		recover()
	}()
	NewRPCStatus(41101, "something is wrong")
}

func TestToRPCCode(t *testing.T) {
	status := []*RPCStatus{
		StatusSuccess,
		StatusCanceled,
		StatusUnknown,
		StatusInvalidParams,
		StatusDeadlineExceeded,
		StatusNotFound,
		StatusAlreadyExists,
		StatusPermissionDenied,
		StatusResourceExhausted,
		StatusFailedPrecondition,
		StatusAborted,
		StatusOutOfRange,
		StatusUnimplemented,
		StatusInternalServerError,
		StatusServiceUnavailable,
		StatusDataLoss,
		StatusUnauthorized,
		StatusTimeout,
		StatusTooManyRequests,
		StatusForbidden,
		StatusLimitExceed,
		StatusMethodNotAllowed,
		StatusAccessDenied,
	}

	var codes []string
	for _, s := range status {
		codes = append(codes, s.ToRPCCode().String())
	}
	t.Log(codes)

	var errors []error
	for i, s := range status {
		if i%2 == 0 {
			errors = append(errors, s.ToRPCErr())
			continue
		}
		errors = append(errors, s.ToRPCErr(s.status.Message()))
	}
	t.Log(errors)

	codeInt := []int{}
	for _, s := range status {
		codeInt = append(codeInt, ToHTTPErr(s.status).code)
	}
	t.Log(codeInt)
}

func TestRCode(t *testing.T) {
	code := RCode(1)
	t.Log("error code is", int(code))

	defer func() {
		recover()
	}()
	code = RCode(101)
}
