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

// Code get code
func (g *RPCStatus) Code() codes.Code {
	return g.status.Code()
}

// Msg get message
func (g *RPCStatus) Msg() string {
	return g.status.Message()
}

// ToRPCErr converted to standard RPC error
func (g *RPCStatus) ToRPCErr(desc ...string) error {
	switch g.status.Code() {
	case StatusCanceled.status.Code():
		return toRPCErr(codes.Canceled, desc...)
	case StatusUnknown.status.Code():
		return toRPCErr(codes.Unknown, desc...)
	case StatusInvalidParams.status.Code():
		return toRPCErr(codes.InvalidArgument, desc...)
	case StatusDeadlineExceeded.status.Code():
		return toRPCErr(codes.DeadlineExceeded, desc...)
	case StatusNotFound.status.Code():
		return toRPCErr(codes.NotFound, desc...)
	case StatusAlreadyExists.status.Code():
		return toRPCErr(codes.AlreadyExists, desc...)
	case StatusPermissionDenied.status.Code():
		return toRPCErr(codes.PermissionDenied, desc...)
	case StatusResourceExhausted.status.Code():
		return toRPCErr(codes.ResourceExhausted, desc...)
	case StatusFailedPrecondition.status.Code():
		return toRPCErr(codes.FailedPrecondition, desc...)
	case StatusAborted.status.Code():
		return toRPCErr(codes.Aborted, desc...)
	case StatusOutOfRange.status.Code():
		return toRPCErr(codes.OutOfRange, desc...)
	case StatusUnimplemented.status.Code():
		return toRPCErr(codes.Unimplemented, desc...)
	case StatusInternalServerError.status.Code():
		return toRPCErr(codes.Internal, desc...)
	case StatusServiceUnavailable.status.Code():
		return toRPCErr(codes.Unavailable, desc...)
	case StatusDataLoss.status.Code():
		return toRPCErr(codes.DataLoss, desc...)
	case StatusUnauthorized.status.Code():
		return toRPCErr(codes.Unauthenticated, desc...)

	case StatusAccessDenied.status.Code():
		return toRPCErr(codes.PermissionDenied, desc...)
	case StatusLimitExceed.status.Code():
		return toRPCErr(codes.ResourceExhausted, desc...)
	case StatusMethodNotAllowed.status.Code():
		return toRPCErr(codes.Unimplemented, desc...)
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
	case StatusCanceled.status.Code():
		return codes.Canceled
	case StatusUnknown.status.Code():
		return codes.Unknown
	case StatusInvalidParams.status.Code():
		return codes.InvalidArgument
	case StatusDeadlineExceeded.status.Code():
		return codes.DeadlineExceeded
	case StatusNotFound.status.Code():
		return codes.NotFound
	case StatusAlreadyExists.status.Code():
		return codes.AlreadyExists
	case StatusPermissionDenied.status.Code():
		return codes.PermissionDenied
	case StatusResourceExhausted.status.Code():
		return codes.ResourceExhausted
	case StatusFailedPrecondition.status.Code():
		return codes.FailedPrecondition
	case StatusAborted.status.Code():
		return codes.Aborted
	case StatusOutOfRange.status.Code():
		return codes.OutOfRange
	case StatusUnimplemented.status.Code():
		return codes.Unimplemented
	case StatusInternalServerError.status.Code():
		return codes.Internal
	case StatusServiceUnavailable.status.Code():
		return codes.Unavailable
	case StatusDataLoss.status.Code():
		return codes.DataLoss
	case StatusUnauthorized.status.Code():
		return codes.Unauthenticated

	case StatusAccessDenied.status.Code():
		return codes.PermissionDenied
	case StatusLimitExceed.status.Code():
		return codes.ResourceExhausted
	case StatusMethodNotAllowed.status.Code():
		return codes.Unimplemented
	}

	return g.status.Code()
}
