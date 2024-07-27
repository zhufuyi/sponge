package errcode

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SkipResponse skip response
var SkipResponse = errors.New("skip response") //nolint

// Responser response interface
type Responser interface {
	Success(ctx *gin.Context, data interface{})
	ParamError(ctx *gin.Context, err error)
	Error(ctx *gin.Context, err error) bool
}

// NewResponser creates a new responser, if isFromRPC=true, it means return from rpc, otherwise default return from http
func NewResponser(isFromRPC bool, httpErrors []*Error, rpcStatus []*RPCStatus) Responser {
	httpErrorsMap := make(map[int]*Error)
	rpcStatusMap := make(map[int]*RPCStatus)

	for _, httpError := range httpErrors {
		if httpError == nil {
			continue
		}
		httpErrorsMap[httpError.Code()] = httpError
	}
	for _, statusError := range rpcStatus {
		if statusError == nil || statusError.status == nil {
			continue
		}
		rpcStatusMap[int(statusError.ToRPCCode())] = statusError
		rpcStatusMap[int(statusError.status.Code())] = statusError
	}

	return &defaultResponse{
		isFromRPC:  isFromRPC,
		httpErrors: httpErrorsMap,
		rpcStatus:  rpcStatusMap,
	}
}

type defaultResponse struct {
	isFromRPC  bool // error comes from grpc, if not, default is from http
	httpErrors map[int]*Error
	rpcStatus  map[int]*RPCStatus
}

func (resp *defaultResponse) response(c *gin.Context, respStatus, code int, msg string, data interface{}) {
	c.JSON(respStatus, map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

// Success response success information
func (resp *defaultResponse) Success(c *gin.Context, data interface{}) {
	resp.response(c, http.StatusOK, 0, "ok", data)
}

// ParamError response parameter error information, does not return an error message
func (resp *defaultResponse) ParamError(c *gin.Context, _ error) {
	resp.response(c, http.StatusOK, InvalidParams.Code(), InvalidParams.Msg(), struct{}{})
}

// Error response error information, if return true, means that the error code is converted to a standard http code,
// otherwise the return http code is always 200
func (resp *defaultResponse) Error(c *gin.Context, err error) bool {
	if resp.isFromRPC {
		// error from rpc and response the corresponding http code
		return resp.handleRPCError(c, err)
	}

	// error from http and response http code
	return resp.handleHTTPError(c, err)
}

// error from grpc
func (resp *defaultResponse) handleRPCError(c *gin.Context, err error) bool {
	st, _ := status.FromError(err)

	// user defined err, response 200
	if st.Code() == codes.Unknown {
		code, msg := parseCodeAndMsg(st.String())
		if code == -1 {
			// non-conforming err
			resp.response(c, http.StatusOK, -1, "unknown error", struct{}{})
		} else {
			// err created using NewRPCStatus
			resp.response(c, http.StatusOK, code, msg, struct{}{})
		}
		return false
	}

	// default error code to http
	switch st.Code() {
	case codes.Internal, StatusInternalServerError.status.Code():
		resp.response(c, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), struct{}{})
		return true
	case codes.Unavailable, StatusServiceUnavailable.status.Code():
		resp.response(c, http.StatusServiceUnavailable, http.StatusServiceUnavailable, http.StatusText(http.StatusServiceUnavailable), struct{}{})
		return true
	}

	// check if you need to return the standard http code
	if strings.Contains(st.Message(), ToHTTPCodeLabel) {
		code := convertToHTTPCode(st.Code())
		msg := strings.ReplaceAll(st.Message(), ToHTTPCodeLabel, "")
		resp.response(c, code, int(st.Code()), msg, struct{}{})
		return true
	}

	// user defined error code to http
	if resp.isUserDefinedRPCErrorCode(c, int(st.Code())) {
		return true
	}

	// response 200
	resp.response(c, http.StatusOK, int(st.Code()), st.Message(), struct{}{})

	return false
}

// error from http
func (resp *defaultResponse) handleHTTPError(c *gin.Context, err error) bool {
	e := ParseError(err)

	// default error code to http
	switch e.Code() {
	case InternalServerError.Code(), http.StatusInternalServerError:
		resp.response(c, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), struct{}{})
		return true
	case ServiceUnavailable.Code(), http.StatusServiceUnavailable:
		resp.response(c, http.StatusServiceUnavailable, http.StatusServiceUnavailable, http.StatusText(http.StatusServiceUnavailable), struct{}{})
		return true
	}

	// user requests to return standard HTTP code, if e.ToHTTPCode() not match, will return of 500
	if e.needHTTPCode {
		msg := strings.ReplaceAll(e.msg, ToHTTPCodeLabel, "")
		resp.response(c, e.ToHTTPCode(), e.code, msg, struct{}{})
		return true
	}

	// user defined error code to http
	if resp.isUserDefinedHTTPErrorCode(c, e.Code()) {
		return true
	}

	// response 200
	resp.response(c, http.StatusOK, e.code, e.msg, struct{}{})
	return false
}

func (resp *defaultResponse) isUserDefinedRPCErrorCode(c *gin.Context, errCode int) bool {
	if v, ok := resp.rpcStatus[errCode]; ok {
		httpCode := ToHTTPErr(v.status).ToHTTPCode()
		msg := http.StatusText(httpCode)
		if msg == "" {
			msg = "unknown error"
		}
		resp.response(c, httpCode, httpCode, msg, struct{}{})
		return true
	}
	return false
}

func (resp *defaultResponse) isUserDefinedHTTPErrorCode(c *gin.Context, errCode int) bool {
	if v, ok := resp.httpErrors[errCode]; ok {
		httpCode := v.ToHTTPCode()
		msg := http.StatusText(httpCode)
		if msg == "" {
			msg = "unknown error"
		}
		resp.response(c, httpCode, httpCode, msg, struct{}{})
		return true
	}
	return false
}

// ToHTTPErr converted to http error
func ToHTTPErr(st *status.Status) *Error { //nolint
	switch st.Code() {
	case StatusSuccess.status.Code(), codes.OK:
		return Success
	case StatusInvalidParams.status.Code(), codes.InvalidArgument:
		return InvalidParams
	case StatusInternalServerError.status.Code(), codes.Internal:
		return InternalServerError
	case StatusUnimplemented.status.Code(), codes.Unimplemented:
		return Unimplemented
	case StatusPermissionDenied.status.Code(), codes.PermissionDenied:
		return PermissionDenied
	}

	switch st.Code() {
	case StatusCanceled.status.Code(), codes.Canceled:
		return Canceled
	case StatusUnknown.status.Code(), codes.Unknown:
		return Unknown
	case StatusDeadlineExceeded.status.Code(), codes.DeadlineExceeded:
		return DeadlineExceeded
	case StatusNotFound.status.Code(), codes.NotFound:
		return NotFound
	case StatusAlreadyExists.status.Code(), codes.AlreadyExists:
		return AlreadyExists
	case StatusResourceExhausted.status.Code(), codes.ResourceExhausted:
		return ResourceExhausted
	case StatusFailedPrecondition.status.Code(), codes.FailedPrecondition:
		return FailedPrecondition
	case StatusAborted.status.Code(), codes.Aborted:
		return Aborted
	case StatusOutOfRange.status.Code(), codes.OutOfRange:
		return OutOfRange
	case StatusServiceUnavailable.status.Code(), codes.Unavailable:
		return ServiceUnavailable
	case StatusDataLoss.status.Code(), codes.DataLoss:
		return DataLoss
	case StatusUnauthorized.status.Code(), codes.Unauthenticated:
		return Unauthorized

	case StatusAccessDenied.status.Code():
		return AccessDenied
	case StatusLimitExceed.status.Code():
		return LimitExceed
	case StatusMethodNotAllowed.status.Code():
		return MethodNotAllowed
	}

	return &Error{
		code: int(st.Code()),
		msg:  st.Message(),
	}
}

func parseCodeAndMsg(errStr string) (int, string) {
	if errStr != "" {
		ss := strings.Split(errStr, "desc = ")
		cm := strings.Split(ss[len(ss)-1], "msg = ")
		if len(cm) == 2 {
			codeStr := strings.ReplaceAll(cm[0], "code = ", "")
			codeStr = strings.ReplaceAll(codeStr, ", ", "")
			code, _ := strconv.Atoi(codeStr)
			msg := cm[1]
			return code, msg
		}
	}
	return -1, errStr
}
