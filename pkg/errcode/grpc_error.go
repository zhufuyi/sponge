package errcode

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCStatus grpc 状态
type GRPCStatus struct {
	status *status.Status
}

var statusCodes = map[codes.Code]string{}

// NewGRPCStatus 新建一个status
func NewGRPCStatus(code codes.Code, msg string) *GRPCStatus {
	if v, ok := statusCodes[code]; ok {
		panic(fmt.Sprintf("grpc status code = %d already exists, please replace with a new error code, old msg = %s", code, v))
	} else {
		statusCodes[code] = msg
	}

	return &GRPCStatus{
		status: status.New(code, msg),
	}
}

// Detail error details
type Detail struct {
	key string
	val interface{}
}

// String detail key-value
func (d *Detail) String() string {
	return fmt.Sprintf("%s: {%v}", d.key, d.val)
}

// Any type key value
func Any(key string, val interface{}) Detail {
	return Detail{
		key: key,
		val: val,
	}
}

// Err return error
func (g *GRPCStatus) Err(details ...Detail) error {
	var dts []string
	for _, detail := range details {
		dts = append(dts, detail.String())
	}
	if len(dts) == 0 {
		return status.Errorf(g.status.Code(), "%s", g.status.Message())
	}
	return status.Errorf(g.status.Code(), "%s details = %v", g.status.Message(), dts)
}

// ToRPCCode 转换为RPC识别的错误码，避免返回Unknown状态码
func ToRPCCode(code codes.Code) codes.Code {
	switch code {
	case StatusInternalServerError.status.Code():
		code = codes.Internal
	case StatusInvalidParams.status.Code():
		code = codes.InvalidArgument
	case StatusUnauthorized.status.Code():
		code = codes.Unauthenticated
	case StatusNotFound.status.Code():
		code = codes.NotFound
	case StatusDeadlineExceeded.status.Code():
		code = codes.DeadlineExceeded
	case StatusAccessDenied.status.Code():
		code = codes.PermissionDenied
	case StatusLimitExceed.status.Code():
		code = codes.ResourceExhausted
	case StatusMethodNotAllowed.status.Code():
		code = codes.Unimplemented
	}

	return code
}
