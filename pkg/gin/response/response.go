package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zhufuyi/sponge/pkg/errcode"

	"github.com/gin-gonic/gin"
)

// Result 输出数据格式
type Result struct {
	Code int         `json:"code"` // 返回码
	Msg  string      `json:"msg"`  // 返回信息说明
	Data interface{} `json:"data"` // 返回数据
}

func newResp(code int, msg string, data interface{}) *Result {
	resp := &Result{
		Code: code,
		Msg:  msg,
	}

	// 保证返回时data字段不为nil，注意resp.Data=[]interface {}时不为nil，经过序列化变成了null
	if data == nil {
		resp.Data = &struct{}{}
	} else {
		resp.Data = data
	}

	return resp
}

var jsonContentType = []string{"application/json; charset=utf-8"}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

func writeJSON(c *gin.Context, code int, res interface{}) {
	c.Writer.WriteHeader(code)
	writeContentType(c.Writer, jsonContentType)
	err := json.NewEncoder(c.Writer).Encode(res)
	if err != nil {
		fmt.Printf("json encode error, err = %s\n", err.Error())
	}
}

func respJSONWithStatusCode(c *gin.Context, code int, msg string, data ...interface{}) {
	var FirstData interface{}
	if len(data) > 0 {
		FirstData = data[0]
	}
	resp := newResp(code, msg, FirstData)

	writeJSON(c, code, resp)
}

// Output 根据http status code返回json数据
func Output(c *gin.Context, code int, msg ...interface{}) {
	switch code {
	case http.StatusOK:
		respJSONWithStatusCode(c, http.StatusOK, "ok", msg...)
	case http.StatusBadRequest:
		respJSONWithStatusCode(c, http.StatusBadRequest, errcode.InvalidParams.Msg(), msg...)
	case http.StatusUnauthorized:
		respJSONWithStatusCode(c, http.StatusUnauthorized, errcode.Unauthorized.Msg(), msg...)
	case http.StatusForbidden:
		respJSONWithStatusCode(c, http.StatusForbidden, errcode.Forbidden.Msg(), msg...)
	case http.StatusNotFound:
		respJSONWithStatusCode(c, http.StatusNotFound, errcode.NotFound.Msg(), msg...)
	case http.StatusRequestTimeout:
		respJSONWithStatusCode(c, http.StatusRequestTimeout, errcode.Timeout.Msg(), msg...)
	case http.StatusConflict:
		respJSONWithStatusCode(c, http.StatusConflict, errcode.AlreadyExists.Msg(), msg...)
	case http.StatusInternalServerError:
		respJSONWithStatusCode(c, http.StatusInternalServerError, errcode.InternalServerError.Msg(), msg...)
	case http.StatusTooManyRequests:
		respJSONWithStatusCode(c, http.StatusTooManyRequests, errcode.LimitExceed.Msg(), msg...)
	case http.StatusServiceUnavailable:
		respJSONWithStatusCode(c, http.StatusServiceUnavailable, errcode.ServiceUnavailable.Msg(), msg...)

	default:
		respJSONWithStatusCode(c, code, http.StatusText(code), msg...)
	}
}

// 状态码统一200，自定义错误码在data.code
func respJSONWith200(c *gin.Context, code int, msg string, data ...interface{}) {
	var FirstData interface{}
	if len(data) > 0 {
		FirstData = data[0]
	}
	resp := newResp(code, msg, FirstData)

	writeJSON(c, http.StatusOK, resp)
}

// Success 正确
func Success(c *gin.Context, data ...interface{}) {
	respJSONWith200(c, 0, "ok", data...)
}

// Error 错误
func Error(c *gin.Context, err *errcode.Error, data ...interface{}) {
	respJSONWith200(c, err.Code(), err.Msg(), data...)
}
