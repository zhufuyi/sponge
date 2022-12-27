package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zhufuyi/sponge/pkg/errcode"

	"github.com/gin-gonic/gin"
)

// Result output data format
type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func newResp(code int, msg string, data interface{}) *Result {
	resp := &Result{
		Code: code,
		Msg:  msg,
	}

	// ensure that the data field is not nil on return, note that it is not nil when resp.data=[]interface {}, it is serialized to null
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

// Output return json data by http status code
// Deprecated: Output use Out() instead.
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

// Out return json data by http status code, converted by errcode
func Out(c *gin.Context, err *errcode.Error, data ...interface{}) {
	code := err.ToHTTPCode()
	switch code {
	case http.StatusOK:
		respJSONWithStatusCode(c, http.StatusOK, "ok", data...)
	case http.StatusBadRequest:
		respJSONWithStatusCode(c, http.StatusBadRequest, err.Msg(), data...)
	case http.StatusUnauthorized:
		respJSONWithStatusCode(c, http.StatusUnauthorized, err.Msg(), data...)
	case http.StatusForbidden:
		respJSONWithStatusCode(c, http.StatusForbidden, err.Msg(), data...)
	case http.StatusNotFound:
		respJSONWithStatusCode(c, http.StatusNotFound, err.Msg(), data...)
	case http.StatusRequestTimeout:
		respJSONWithStatusCode(c, http.StatusRequestTimeout, err.Msg(), data...)
	case http.StatusConflict:
		respJSONWithStatusCode(c, http.StatusConflict, err.Msg(), data...)
	case http.StatusInternalServerError:
		respJSONWithStatusCode(c, http.StatusInternalServerError, err.Msg(), data...)
	case http.StatusTooManyRequests:
		respJSONWithStatusCode(c, http.StatusTooManyRequests, err.Msg(), data...)
	case http.StatusServiceUnavailable:
		respJSONWithStatusCode(c, http.StatusServiceUnavailable, err.Msg(), data...)

	default:
		respJSONWithStatusCode(c, http.StatusNotExtended, err.Msg(), data...)
	}
}

// status code flat 200, custom error codes in data.code
func respJSONWith200(c *gin.Context, code int, msg string, data ...interface{}) {
	var FirstData interface{}
	if len(data) > 0 {
		FirstData = data[0]
	}
	resp := newResp(code, msg, FirstData)

	writeJSON(c, http.StatusOK, resp)
}

// Success return success
func Success(c *gin.Context, data ...interface{}) {
	respJSONWith200(c, 0, "ok", data...)
}

// Error return error
func Error(c *gin.Context, err *errcode.Error, data ...interface{}) {
	respJSONWith200(c, err.Code(), err.Msg(), data...)
}
