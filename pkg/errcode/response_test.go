package errcode

import (
	"errors"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zhufuyi/sponge/pkg/utils"
)

func runHTTPServer(isFromRPC bool) string {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	httpErrors := []*Error{Forbidden, TooManyRequests, MethodNotAllowed}
	rpcStatus := []*RPCStatus{StatusDeadlineExceeded, StatusPermissionDenied, StatusAlreadyExists}
	resp := NewResponser(isFromRPC, httpErrors, rpcStatus)

	r.GET("/ping", func(c *gin.Context) {
		resp.Success(c, "ping")
	})
	r.GET("/unknown", func(c *gin.Context) {
		resp.Error(c, errors.New("unknown error"))
	})

	if isFromRPC {
		r.GET("/grpc/params", func(c *gin.Context) {
			resp.Error(c, StatusInvalidParams.Err("email is required"))
		})
		r.GET("/grpc/src_params", func(c *gin.Context) {
			resp.Error(c, StatusInvalidParams.ToRPCErr("name is required"))
		})
		r.GET("/grpc/mark_params", func(c *gin.Context) {
			resp.Error(c, StatusInvalidParams.ErrToHTTP("id must be a number"))
		})

		r.GET("/grpc/internal", func(c *gin.Context) {
			resp.Error(c, StatusInternalServerError.Err())
		})
		r.GET("/grpc/src_internal", func(c *gin.Context) {
			resp.Error(c, StatusInternalServerError.ToRPCErr())
		})
		r.GET("/grpc/mark_internal", func(c *gin.Context) {
			resp.Error(c, StatusInternalServerError.ErrToHTTP())
		})

		r.GET("/grpc/unavailable", func(c *gin.Context) {
			resp.Error(c, StatusServiceUnavailable.Err())
		})
		r.GET("/grpc/src_unavailable", func(c *gin.Context) {
			resp.Error(c, StatusServiceUnavailable.ToRPCErr())
		})
		r.GET("/grpc/mark_unavailable", func(c *gin.Context) {
			resp.Error(c, StatusServiceUnavailable.ErrToHTTP())
		})

		r.GET("/grpc/notfound", func(c *gin.Context) {
			resp.Error(c, StatusNotFound.Err())
		})
		r.GET("/grpc/src_notfound", func(c *gin.Context) {
			resp.Error(c, StatusNotFound.ToRPCErr())
		})
		r.GET("/grpc/mark_notfound", func(c *gin.Context) {
			resp.Error(c, StatusNotFound.ErrToHTTP())
		})

		r.GET("/grpc/permission", func(c *gin.Context) {
			resp.Error(c, StatusPermissionDenied.Err())
		})
		r.GET("/grpc/src_permission", func(c *gin.Context) {
			resp.Error(c, StatusPermissionDenied.ToRPCErr())
		})
		r.GET("/grpc/mark_permission", func(c *gin.Context) {
			resp.Error(c, StatusPermissionDenied.ErrToHTTP())
		})

		r.GET("/grpc/conflict", func(c *gin.Context) {
			resp.Error(c, StatusConflict.Err())
		})
		r.GET("/grpc/src_conflict", func(c *gin.Context) {
			resp.Error(c, StatusConflict.ToRPCErr())
		})
		r.GET("/grpc/mark_conflict", func(c *gin.Context) {
			resp.Error(c, StatusConflict.ErrToHTTP())
		})
	} else {
		r.GET("/http/params", func(c *gin.Context) {
			resp.Error(c, InvalidParams.Err("name is required"))
		})
		r.GET("/http/mark_params", func(c *gin.Context) {
			resp.Error(c, InvalidParams.ErrToHTTP("id must be a number"))
		})

		r.GET("/http/internal", func(c *gin.Context) {
			resp.Error(c, InternalServerError.Err())
		})
		r.GET("/http/mark_internal", func(c *gin.Context) {
			resp.Error(c, InternalServerError.ErrToHTTP())
		})

		r.GET("/http/notfound", func(c *gin.Context) {
			resp.Error(c, NotFound.Err())
		})
		r.GET("/http/mark_notfound", func(c *gin.Context) {
			resp.Error(c, NotFound.ErrToHTTP())
		})

		r.GET("/http/forbidden", func(c *gin.Context) {
			resp.Error(c, Forbidden.Err())
		})
		r.GET("/http/mark_forbidden", func(c *gin.Context) {
			resp.Error(c, Forbidden.ErrToHTTP())
		})

		r.GET("/http/too_many", func(c *gin.Context) {
			resp.Error(c, TooManyRequests.Err())
		})
		r.GET("/http/mark_too_many", func(c *gin.Context) {
			resp.Error(c, TooManyRequests.ErrToHTTP())
		})

		r.GET("/http/conflict", func(c *gin.Context) {
			resp.Error(c, Conflict.Err())
		})
		r.GET("/http/mark_conflict", func(c *gin.Context) {
			resp.Error(c, Conflict.ErrToHTTP())
		})
	}

	go func() {
		err := r.Run(serverAddr)
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Millisecond * 200)

	return requestAddr
}

func TestRPCResponse(t *testing.T) {
	requestAddr := runHTTPServer(true)

	result, err := http.Get(requestAddr + "/ping")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusOK)

	result, err = http.Get(requestAddr + "/unknown")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusOK)

	result, err = http.Get(requestAddr + "/grpc/params")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusOK)
	data, e := io.ReadAll(result.Body)
	t.Log(string(data), e)
	result, err = http.Get(requestAddr + "/grpc/src_params")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusOK)
	data, e = io.ReadAll(result.Body)
	t.Log(string(data), e)
	result, err = http.Get(requestAddr + "/grpc/mark_params")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusBadRequest)
	data, e = io.ReadAll(result.Body)
	t.Log(string(data), e)

	result, err = http.Get(requestAddr + "/grpc/internal")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusInternalServerError)
	result, err = http.Get(requestAddr + "/grpc/src_internal")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusInternalServerError)
	result, err = http.Get(requestAddr + "/grpc/mark_internal")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

	result, err = http.Get(requestAddr + "/grpc/unavailable")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusServiceUnavailable)
	result, err = http.Get(requestAddr + "/grpc/src_unavailable")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusServiceUnavailable)
	result, err = http.Get(requestAddr + "/grpc/mark_unavailable")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusServiceUnavailable)

	result, err = http.Get(requestAddr + "/grpc/notfound")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusOK)
	result, err = http.Get(requestAddr + "/grpc/src_notfound")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusOK)
	result, err = http.Get(requestAddr + "/grpc/mark_notfound")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusNotFound)

	result, err = http.Get(requestAddr + "/grpc/permission")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusUnauthorized)
	result, err = http.Get(requestAddr + "/grpc/src_permission")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusUnauthorized)
	result, err = http.Get(requestAddr + "/grpc/mark_permission")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusUnauthorized)

	result, err = http.Get(requestAddr + "/grpc/conflict")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusOK)
	result, err = http.Get(requestAddr + "/grpc/src_conflict")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusConflict)
	result, err = http.Get(requestAddr + "/grpc/mark_conflict")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusConflict)
}

func TestHTTPResponse(t *testing.T) {
	requestAddr := runHTTPServer(false)

	result, err := http.Get(requestAddr + "/ping")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusOK)

	result, err = http.Get(requestAddr + "/unknown")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusOK)

	result, err = http.Get(requestAddr + "/http/params")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusOK)
	data, e := io.ReadAll(result.Body)
	t.Log(string(data), e)
	result, err = http.Get(requestAddr + "/http/mark_params")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusBadRequest)
	data, e = io.ReadAll(result.Body)
	t.Log(string(data), e)

	result, err = http.Get(requestAddr + "/http/internal")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusInternalServerError)
	result, err = http.Get(requestAddr + "/http/mark_internal")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

	result, err = http.Get(requestAddr + "/http/notfound")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusOK)
	result, err = http.Get(requestAddr + "/http/mark_notfound")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusNotFound)

	result, err = http.Get(requestAddr + "/http/forbidden")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusForbidden)
	result, err = http.Get(requestAddr + "/http/mark_forbidden")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusForbidden)

	result, err = http.Get(requestAddr + "/http/too_many")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusTooManyRequests)
	result, err = http.Get(requestAddr + "/http/mark_too_many")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusTooManyRequests)

	result, err = http.Get(requestAddr + "/http/conflict")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusOK)
	result, err = http.Get(requestAddr + "/http/mark_conflict")
	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusConflict)
}

func TestParseCodeAndMsgError(t *testing.T) {
	errStr := "rpc error: code = Unknown desc = rpc error: code = Unknown desc = code = 204011, msg = wrong account or password"
	err := errors.New(errStr)
	st, _ := status.FromError(err)
	if st.Code() == codes.Unknown {
		code, msg := parseCodeAndMsg(st.String())
		t.Log("regexp: ", code, msg)
	}
	code, msg := parseCodeAndMsg2(st.String())
	t.Log("strings: ", code, msg)

	st, _ = status.FromError(SkipResponse)
	t.Log(st.Code(), st.Message())
}

var mcReg = regexp.MustCompile(`code\s*=\s*(\d+),\s*msg\s*=\s*(.+)`)

func parseCodeAndMsg2(errStr string) (int, string) {
	matches := mcReg.FindStringSubmatch(errStr)
	if len(matches) == 3 {
		code, _ := strconv.Atoi(matches[1])
		msg := matches[2]
		return code, msg
	}
	return 0, errStr
}

func BenchmarkName(b *testing.B) {
	errStr := "rpc error: code = Unknown desc = rpc error: code = Unknown desc = code = 204011, msg = wrong account or password"
	err := errors.New(errStr)
	st, _ := status.FromError(err)
	if st.Code() == codes.Unknown {
		for i := 0; i < b.N; i++ {
			parseCodeAndMsg(st.String())
			//parseCodeAndMsg2(st.String())
		}
	}
}
