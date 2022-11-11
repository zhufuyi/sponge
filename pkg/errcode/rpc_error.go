package errcode

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RPCStatus rpc status
type RPCStatus struct {
	status *status.Status
}

var statusCodes = map[codes.Code]string{}

// NewRPCStatus create a new rpc status
func NewRPCStatus(code codes.Code, msg string) *RPCStatus {
	if v, ok := statusCodes[code]; ok {
		panic(fmt.Sprintf("grpc status code = %d already exists, please replace with a new error code, old msg = %s", code, v))
	} else {
		statusCodes[code] = msg
	}

	return &RPCStatus{
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
func (g *RPCStatus) Err(details ...Detail) error {
	var dts []string
	for _, detail := range details {
		dts = append(dts, detail.String())
	}
	if len(dts) == 0 {
		return status.Errorf(g.status.Code(), "%s", g.status.Message())
	}
	return status.Errorf(g.status.Code(), "%s details = %s", g.status.Message(), dts)
}

// ToRPCErr converted to standard RPC error
func (g *RPCStatus) ToRPCErr(desc ...string) error {
	switch g.status.Code() {
	case StatusInternalServerError.status.Code():
		return toRPCErr(codes.Internal, desc...)
	case StatusInvalidParams.status.Code():
		return toRPCErr(codes.InvalidArgument, desc...)
	case StatusUnauthorized.status.Code():
		return toRPCErr(codes.Unauthenticated, desc...)
	case StatusNotFound.status.Code():
		return toRPCErr(codes.NotFound, desc...)
	case StatusDeadlineExceeded.status.Code():
		return toRPCErr(codes.DeadlineExceeded, desc...)
	case StatusAccessDenied.status.Code():
		return toRPCErr(codes.PermissionDenied, desc...)
	case StatusLimitExceed.status.Code():
		return toRPCErr(codes.ResourceExhausted, desc...)
	case StatusMethodNotAllowed.status.Code():
		return toRPCErr(codes.Unimplemented, desc...)
	case StatusServiceUnavailable.status.Code():
		return toRPCErr(codes.Unavailable, desc...)
	}

	return g.status.Err()
}

func toRPCErr(code codes.Code, descs ...string) error {
	var desc string
	if len(descs) > 0 {
		desc = strings.Join(descs, ", ")
	} else {
		desc = code.String()
	}
	return status.New(code, desc).Err()
}

// ToRPCCode converted to standard RPC error code
func (g *RPCStatus) ToRPCCode() codes.Code {
	switch g.status.Code() {
	case StatusInternalServerError.status.Code():
		return codes.Internal
	case StatusInvalidParams.status.Code():
		return codes.InvalidArgument
	case StatusUnauthorized.status.Code():
		return codes.Unauthenticated
	case StatusNotFound.status.Code():
		return codes.NotFound
	case StatusDeadlineExceeded.status.Code():
		return codes.DeadlineExceeded
	case StatusAccessDenied.status.Code():
		return codes.PermissionDenied
	case StatusLimitExceed.status.Code():
		return codes.ResourceExhausted
	case StatusMethodNotAllowed.status.Code():
		return codes.Unimplemented
	case StatusServiceUnavailable.status.Code():
		return codes.Unavailable
	}

	return g.status.Code()
}
