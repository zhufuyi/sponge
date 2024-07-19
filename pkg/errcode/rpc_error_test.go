package errcode

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/utils"
)

func TestRPCStatus(t *testing.T) {
	st := NewRPCStatus(41101, "something is wrong")
	err := st.Err()
	assert.Error(t, err)
	err = st.Err(Any("foo", "bar"))
	assert.Error(t, err)
	code := st.Code()
	assert.Equal(t, int(code), 41101)
	msg := st.Msg()
	assert.Equal(t, msg, "something is wrong")

	defer func() {
		if e := recover(); e != nil {
			t.Log(e)
		}
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

func TestConvertToHTTPCode(t *testing.T) {
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

	var codes []int
	for _, s := range status {
		codes = append(codes, convertToHTTPCode(s.Code()))
	}
	t.Log(codes)
}

func TestRCode(t *testing.T) {
	code := RCode(1)
	t.Log("error code is", int(code))

	defer func() {
		recover()
	}()
	code = RCode(101)
}

func TestHandlers(t *testing.T) {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/codes", gin.WrapF(ListGRPCErrCodes))
	r.GET("/config", gin.WrapF(ShowConfig([]byte(`{"foo": "bar"}`))))

	go func() {
		_ = r.Run(serverAddr)
	}()

	time.Sleep(time.Millisecond * 200)
	resp, err := http.Get(requestAddr + "/codes")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	resp, err = http.Get(requestAddr + "/config")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	time.Sleep(time.Second)
}
