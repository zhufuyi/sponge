package errcode

// rpc system level error code with status prefix, error code range 30000~40000
var (
	StatusSuccess = NewRPCStatus(0, "ok")

	StatusCanceled            = NewRPCStatus(300001, "Canceled")
	StatusUnknown             = NewRPCStatus(300002, "Unknown")
	StatusInvalidParams       = NewRPCStatus(300003, "Invalid Parameter")
	StatusDeadlineExceeded    = NewRPCStatus(300004, "Deadline Exceeded")
	StatusNotFound            = NewRPCStatus(300005, "Not Found")
	StatusAlreadyExists       = NewRPCStatus(300006, "Already Exists")
	StatusPermissionDenied    = NewRPCStatus(300007, "Permission Denied")
	StatusResourceExhausted   = NewRPCStatus(300008, "Resource Exhausted")
	StatusFailedPrecondition  = NewRPCStatus(300009, "Failed Precondition")
	StatusAborted             = NewRPCStatus(300010, "Aborted")
	StatusOutOfRange          = NewRPCStatus(300011, "Out Of Range")
	StatusUnimplemented       = NewRPCStatus(300012, "Unimplemented")
	StatusInternalServerError = NewRPCStatus(300013, "Internal Server Error")
	StatusServiceUnavailable  = NewRPCStatus(300014, "Service Unavailable")
	StatusDataLoss            = NewRPCStatus(300015, "Data Loss")
	StatusUnauthorized        = NewRPCStatus(300016, "Unauthorized")

	StatusTimeout          = NewRPCStatus(300017, "Request Timeout")
	StatusTooManyRequests  = NewRPCStatus(300018, "Too Many Requests")
	StatusForbidden        = NewRPCStatus(300019, "Forbidden")
	StatusLimitExceed      = NewRPCStatus(300020, "Limit Exceed")
	StatusMethodNotAllowed = NewRPCStatus(300021, "Method Not Allowed")
	StatusAccessDenied     = NewRPCStatus(300022, "Access Denied")
	StatusConflict         = NewRPCStatus(300023, "Conflict")
)
