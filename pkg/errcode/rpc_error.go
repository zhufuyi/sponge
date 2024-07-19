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
	return fmt.Sprintf("%s: %v", d.key, d.val)
}

// Any type key value
func Any(key string, val interface{}) Detail {
	return Detail{
		key: key,
		val: val,
	}
}

// Err return error
func (s *RPCStatus) Err(details ...Detail) error {
	var dts []string
	for _, detail := range details {
		dts = append(dts, detail.String())
	}
	if len(dts) == 0 {
		return status.Errorf(s.status.Code(), "%s", s.status.Message())
	}
	return status.Errorf(s.status.Code(), "%s details = %s", s.status.Message(), dts)
}

// ErrToHTTP convert to standard error add ToHTTPCodeLabel to error message
func (s *RPCStatus) ErrToHTTP() error {
	return status.Errorf(s.status.Code(), "%s %s", s.status.Message(), ToHTTPCodeLabel)
}

// Code get code
func (s *RPCStatus) Code() codes.Code {
	return s.status.Code()
}

// Msg get message
func (s *RPCStatus) Msg() string {
	return s.status.Message()
}

// ToRPCErr converted to standard RPC error
func (s *RPCStatus) ToRPCErr(desc ...string) error {
	switch s.status.Code() {
	case StatusInvalidParams.status.Code():
		return toRPCErr(codes.InvalidArgument, desc...)
	case StatusInternalServerError.status.Code():
		return toRPCErr(codes.Internal, desc...)
	}

	switch s.status.Code() {
	case StatusCanceled.status.Code():
		return toRPCErr(codes.Canceled, desc...)
	case StatusUnknown.status.Code():
		return toRPCErr(codes.Unknown, desc...)
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

	return s.status.Err()
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
func (s *RPCStatus) ToRPCCode() codes.Code {
	switch s.status.Code() {
	case StatusInvalidParams.status.Code():
		return codes.InvalidArgument
	case StatusInternalServerError.status.Code():
		return codes.Internal
	case StatusUnimplemented.status.Code():
		return codes.Unimplemented
	}

	switch s.status.Code() {
	case StatusPermissionDenied.status.Code():
		return codes.PermissionDenied
	case StatusCanceled.status.Code():
		return codes.Canceled
	case StatusUnknown.status.Code():
		return codes.Unknown
	case StatusDeadlineExceeded.status.Code():
		return codes.DeadlineExceeded
	case StatusNotFound.status.Code():
		return codes.NotFound
	case StatusAlreadyExists.status.Code():
		return codes.AlreadyExists
	case StatusResourceExhausted.status.Code():
		return codes.ResourceExhausted
	case StatusFailedPrecondition.status.Code():
		return codes.FailedPrecondition
	case StatusAborted.status.Code():
		return codes.Aborted
	case StatusOutOfRange.status.Code():
		return codes.OutOfRange
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

	return s.status.Code()
}

// converted grpc code to http code
func convertToHTTPCode(code codes.Code) int {
	switch code {
	case StatusSuccess.status.Code():
		return http.StatusOK
	case codes.InvalidArgument, StatusInvalidParams.status.Code():
		return http.StatusBadRequest
	case codes.Internal, StatusInternalServerError.status.Code():
		return http.StatusInternalServerError
	case codes.Unimplemented, StatusUnimplemented.status.Code():
		return http.StatusNotImplemented
	case codes.NotFound, StatusNotFound.status.Code():
		return http.StatusNotFound
	case StatusForbidden.status.Code(), StatusAccessDenied.status.Code():
		return http.StatusForbidden
	}

	switch code {
	case StatusTimeout.status.Code():
		return http.StatusRequestTimeout
	case StatusTooManyRequests.status.Code(), StatusLimitExceed.status.Code():
		return http.StatusTooManyRequests
	case codes.FailedPrecondition, StatusFailedPrecondition.status.Code():
		return http.StatusPreconditionFailed
	case codes.Unavailable, StatusServiceUnavailable.status.Code():
		return http.StatusServiceUnavailable
	case codes.Unauthenticated, StatusUnauthorized.status.Code():
		return http.StatusUnauthorized
	case codes.PermissionDenied, StatusPermissionDenied.status.Code():
		return http.StatusUnauthorized
	case StatusLimitExceed.status.Code():
		return http.StatusTooManyRequests
	case StatusMethodNotAllowed.status.Code():
		return http.StatusMethodNotAllowed
	}

	return http.StatusInternalServerError
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
