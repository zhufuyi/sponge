package errcode

import (
	"errors"
	"fmt"
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
	r.GET("/err", func(c *gin.Context) {
		isIgnore := resp.Error(c, errors.New("unknown error"))
		fmt.Println("/err", isIgnore)
	})

	if isFromRPC {
		r.GET("/err1", func(c *gin.Context) {
			isIgnore := resp.Error(c, StatusServiceUnavailable.ToRPCErr())
			fmt.Println("/err1", isIgnore)
		})
		r.GET("/err2", func(c *gin.Context) {
			isIgnore := resp.Error(c, StatusInternalServerError.ToRPCErr())
			fmt.Println("/err2", isIgnore)
		})
		r.GET("/err3", func(c *gin.Context) {
			isIgnore := resp.Error(c, StatusNotFound.Err())
			fmt.Println("/err3", isIgnore)
		})
		r.GET("/rpc/userDefine/err1", func(c *gin.Context) {
			isIgnore := resp.Error(c, StatusDeadlineExceeded.Err())
			fmt.Println("/rpc/userDefine/err1", isIgnore)
		})
		r.GET("/rpc/userDefine/err2", func(c *gin.Context) {
			isIgnore := resp.Error(c, StatusPermissionDenied.Err())
			fmt.Println("/rpc/userDefine/err2", isIgnore)
		})
	} else {
		r.GET("/err4", func(c *gin.Context) {
			isIgnore := resp.Error(c, InternalServerError.Err())
			fmt.Println("/err4", isIgnore)
		})
		r.GET("/err5", func(c *gin.Context) {
			isIgnore := resp.Error(c, NotFound.Err())
			fmt.Println("/err5", isIgnore)
		})
		r.GET("/http/userDefine/err1", func(c *gin.Context) {
			isIgnore := resp.Error(c, Forbidden.Err())
			fmt.Println("/http/userDefine/err1", isIgnore)
		})
		r.GET("/http/userDefine/err2", func(c *gin.Context) {
			isIgnore := resp.Error(c, TooManyRequests.Err())
			fmt.Println("/http/userDefine/err2", isIgnore)
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

func TestNewResponse(t *testing.T) {
	resp := NewResponse(true)
	assert.NotNil(t, resp)
}

func TestRPCResponse(t *testing.T) {
	requestAddr := runHTTPServer(true)

	result, err := http.Get(requestAddr + "/ping")
	assert.NoError(t, err)
	t.Log(result.StatusCode)

	result, err = http.Get(requestAddr + "/err")
	assert.NoError(t, err)
	t.Log(result.StatusCode)

	result, err = http.Get(requestAddr + "/err1")
	assert.NoError(t, err)
	t.Log(result.StatusCode)

	result, err = http.Get(requestAddr + "/err2")
	assert.NoError(t, err)
	t.Log(result.StatusCode)

	result, err = http.Get(requestAddr + "/err3")
	assert.NoError(t, err)
	t.Log(result.StatusCode)

	result, err = http.Get(requestAddr + "/rpc/userDefine/err1")
	assert.NoError(t, err)
	t.Log(result.StatusCode)

	result, err = http.Get(requestAddr + "/rpc/userDefine/err2")
	assert.NoError(t, err)
	t.Log(result.StatusCode)
}

func TestHTTPResponse(t *testing.T) {
	requestAddr := runHTTPServer(false)

	result, err := http.Get(requestAddr + "/ping")
	assert.NoError(t, err)
	t.Log(result.StatusCode)

	result, err = http.Get(requestAddr + "/err")
	assert.NoError(t, err)
	t.Log(result.StatusCode)

	result, err = http.Get(requestAddr + "/err4")
	assert.NoError(t, err)
	t.Log(result.StatusCode)

	result, err = http.Get(requestAddr + "/err5")
	assert.NoError(t, err)
	t.Log(result.StatusCode)

	result, err = http.Get(requestAddr + "/http/userDefine/err1")
	assert.NoError(t, err)
	t.Log(result.StatusCode)

	result, err = http.Get(requestAddr + "/http/userDefine/err2")
	assert.NoError(t, err)
	t.Log(result.StatusCode)
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
