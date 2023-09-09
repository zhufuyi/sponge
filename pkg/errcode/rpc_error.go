package errcode

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var grpcErrCodes = map[int]string{}

// RPCStatus rpc status
type RPCStatus struct {
	status *status.Status
}

var statusCodes = map[codes.Code]string{}

// NewRPCStatus create a new rpc status
func NewRPCStatus(code codes.Code, msg string) *RPCStatus {
	if v, ok := statusCodes[code]; ok {
		panic(fmt.Sprintf(`grpc status code = %d already exists, please define a new error code,
msg1 = %s
msg2 = %s
`, code, v, msg))
	}

	grpcErrCodes[int(code)] = msg
	statusCodes[code] = msg
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
	case StatusInvalidParams.status.Code():
		return codes.InvalidArgument
	case StatusInternalServerError.status.Code():
		return codes.Internal
	case StatusUnimplemented.status.Code():
		return codes.Unimplemented
	case StatusPermissionDenied.status.Code():
		return codes.PermissionDenied
	}

	switch g.status.Code() {
	case StatusCanceled.status.Code():
		return codes.Canceled
	case StatusUnknown.status.Code():
		return codes.Unknown
	//case StatusInvalidParams.status.Code():
	//	return codes.InvalidArgument
	case StatusDeadlineExceeded.status.Code():
		return codes.DeadlineExceeded
	case StatusNotFound.status.Code():
		return codes.NotFound
	case StatusAlreadyExists.status.Code():
		return codes.AlreadyExists
	//case StatusPermissionDenied.status.Code():
	//	return codes.PermissionDenied
	case StatusResourceExhausted.status.Code():
		return codes.ResourceExhausted
	case StatusFailedPrecondition.status.Code():
		return codes.FailedPrecondition
	case StatusAborted.status.Code():
		return codes.Aborted
	case StatusOutOfRange.status.Code():
		return codes.OutOfRange
	//case StatusUnimplemented.status.Code():
	//	return codes.Unimplemented
	//case StatusInternalServerError.status.Code():
	//	return codes.Internal
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

// ErrInfo error info
type ErrInfo struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func getErrorInfo(codeInfo map[int]string) []ErrInfo {
	var keys []int
	for key := range codeInfo {
		keys = append(keys, key)
	}

	sort.Ints(keys)
	eis := []ErrInfo{}
	for _, key := range keys {
		eis = append(eis, ErrInfo{
			Code: key,
			Msg:  codeInfo[key],
		})
	}

	return eis
}

// ListGRPCErrCodes list grpc error codes, http handle func
func ListGRPCErrCodes(w http.ResponseWriter, _ *http.Request) {
	eis := getErrorInfo(grpcErrCodes)

	jsonData, err := json.Marshal(&eis)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ShowConfig show config info
// @Summary show config info
// @Description show config info
// @Tags system
// @Accept  json
// @Produce  json
// @Router /config [get]
func ShowConfig(jsonData []byte) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
