package errcode

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"

	"github.com/zhufuyi/sponge/pkg/utils"
)

var rpcStatus = []*RPCStatus{
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
	StatusConflict,
}

func TestRPCStatus(t *testing.T) {
	st := NewRPCStatus(401101, "something is wrong")
	err := st.Err()
	assert.Error(t, err)
	err = st.Err("another thing is wrong")

	s, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, s.Code(), st.Code())
	assert.Equal(t, s.Message(), "another thing is wrong")

	code := st.Code()
	assert.Equal(t, int(code), 401101)
	msg := st.Msg()
	assert.Equal(t, msg, "something is wrong")

	defer func() {
		if e := recover(); e != nil {
			t.Log(e)
		}
	}()
	NewRPCStatus(401101, "something is wrong 2")
}

func TestToRPCCode(t *testing.T) {
	var codes []string
	for _, s := range rpcStatus {
		codes = append(codes, s.ToRPCCode().String())
	}
	t.Log(codes)

	var errors []error
	for i, s := range rpcStatus {
		if i%2 == 0 {
			errors = append(errors, s.ToRPCErr())
			continue
		}
		errors = append(errors, s.ToRPCErr(s.status.Message()))
	}
	t.Log(errors)

	codeInt := []int{}
	for _, s := range rpcStatus {
		codeInt = append(codeInt, ToHTTPErr(s.status).code)
	}
	t.Log(codeInt)
}

func TestConvertToHTTPCode(t *testing.T) {
	var codes []int
	for _, s := range rpcStatus {
		codes = append(codes, convertToHTTPCode(s.Code()))
	}
	t.Log(codes)

}

func TestGetStatusCode(t *testing.T) {
	t.Log(GetStatusCode(fmt.Errorf("reason for error")))

	for _, s := range rpcStatus {
		t.Log(s.Code(), "|",
			GetStatusCode(s.Err()),
			GetStatusCode(s.Err("reason for error")), "|",

			GetStatusCode(s.ToRPCErr()),
			GetStatusCode(s.ToRPCErr("reason for error")), "|",

			GetStatusCode(s.ErrToHTTP()),
			GetStatusCode(s.ErrToHTTP("reason for error")),
		)
	}
}

func TestRCode(t *testing.T) {
	code := RCode(1)
	t.Log("error code is", int(code))

	defer func() {
		recover()
	}()
	code = RCode(1001)
	t.Log("error code is", int(code))
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
