package errcode

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Responser response interface
type Responser interface {
	Success(ctx *gin.Context, data interface{})
	ParamError(ctx *gin.Context, err error)
	Error(ctx *gin.Context, err error) bool
}

// Deprecated: NewResponse use NewResponser instead
func NewResponse(isFromRPC bool) Responser {
	return &defaultResponse{isFromRPC, make(map[int]*Error), make(map[int]*RPCStatus)}
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

func (resp *defaultResponse) response(c *gin.Context, status, code int, msg string, data interface{}) {
	c.JSON(status, map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

// Success response success information
func (resp *defaultResponse) Success(c *gin.Context, data interface{}) {
	resp.response(c, http.StatusOK, 0, "ok", data)
}

// ParamError response parameter error information
func (resp *defaultResponse) ParamError(c *gin.Context, err error) {
	resp.response(c, http.StatusOK, InvalidParams.Code(), InvalidParams.Msg(), struct{}{})
}

// Error response error information, if return true, means that the error code is converted to a standard http code,
// otherwise the return code is always 200
func (resp *defaultResponse) Error(c *gin.Context, err error) bool {
	_ = c.Error(err)

	// error from rpc
	if resp.isFromRPC {
		st, ok := status.FromError(err)
		if !ok {
			resp.response(c, http.StatusOK, -1, "unknown error", struct{}{})
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

		// User defined error code to http
		if resp.isUserDefinedRPCErrorCode(c, int(st.Code())) {
			return true
		}

		e := ToHTTPErr(st)
		resp.response(c, http.StatusOK, e.code, e.msg, struct{}{})
		return false
	}

	// error from http
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

	// User defined error code to http
	if resp.isUserDefinedHTTPErrorCode(c, e.Code()) {
		return true
	}

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
func ToHTTPErr(st *status.Status) *Error {
	switch st.Code() {
	case StatusSuccess.status.Code(), codes.OK:
		return Success
	case StatusCanceled.status.Code(), codes.Canceled:
		return Canceled
	case StatusUnknown.status.Code(), codes.Unknown:
		return Unknown
	case StatusInvalidParams.status.Code(), codes.InvalidArgument:
		return InvalidParams
	case StatusDeadlineExceeded.status.Code(), codes.DeadlineExceeded:
		return DeadlineExceeded
	case StatusNotFound.status.Code(), codes.NotFound:
		return NotFound
	case StatusAlreadyExists.status.Code(), codes.AlreadyExists:
		return AlreadyExists
	case StatusPermissionDenied.status.Code(), codes.PermissionDenied:
		return PermissionDenied
	case StatusResourceExhausted.status.Code(), codes.ResourceExhausted:
		return ResourceExhausted
	case StatusFailedPrecondition.status.Code(), codes.FailedPrecondition:
		return FailedPrecondition
	case StatusAborted.status.Code(), codes.Aborted:
		return Aborted
	case StatusOutOfRange.status.Code(), codes.OutOfRange:
		return OutOfRange
	case StatusUnimplemented.status.Code(), codes.Unimplemented:
		return Unimplemented
	case StatusInternalServerError.status.Code(), codes.Internal:
		return InternalServerError
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
		code: getCodeInt(st),
		msg:  st.Message(),
	}
}

func getCodeInt(st *status.Status) int {
	code := st.Code().String()
	if len(code) <= 6 {
		return -1
	}

	codeStr := code[5 : len(code)-1]
	codeInt, err := strconv.Atoi(codeStr)
	if err != nil {
		return -1
	}
	return codeInt
}
