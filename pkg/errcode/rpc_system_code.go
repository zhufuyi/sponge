package errcode

// rpc system level error code with status prefix, error code range 30000~40000
var (
	StatusSuccess = NewRPCStatus(0, "ok")

	StatusCanceled            = NewRPCStatus(30001, "Canceled")
	StatusUnknown             = NewRPCStatus(30002, "Unknown")
	StatusInvalidParams       = NewRPCStatus(30003, "Invalid Parameter")
	StatusDeadlineExceeded    = NewRPCStatus(3004, "Deadline Exceeded")
	StatusNotFound            = NewRPCStatus(30005, "Not Found")
	StatusAlreadyExists       = NewRPCStatus(30006, "Already Exists")
	StatusPermissionDenied    = NewRPCStatus(30007, "Permission Denied")
	StatusResourceExhausted   = NewRPCStatus(30008, "Resource Exhausted")
	StatusFailedPrecondition  = NewRPCStatus(30009, "Failed Precondition")
	StatusAborted             = NewRPCStatus(30010, "Aborted")
	StatusOutOfRange          = NewRPCStatus(30011, "Out Of Range")
	StatusUnimplemented       = NewRPCStatus(30012, "Unimplemented")
	StatusInternalServerError = NewRPCStatus(30013, "Internal Server Error")
	StatusServiceUnavailable  = NewRPCStatus(30014, "Service Unavailable")
	StatusDataLoss            = NewRPCStatus(300115, "Data Loss")
	StatusUnauthorized        = NewRPCStatus(30016, "Unauthorized")

	StatusTimeout          = NewRPCStatus(30017, "Request Timeout")
	StatusTooManyRequests  = NewRPCStatus(30018, "Too Many Requests")
	StatusForbidden        = NewRPCStatus(30019, "Forbidden")
	StatusLimitExceed      = NewRPCStatus(30020, "Limit Exceed")
	StatusMethodNotAllowed = NewRPCStatus(30021, "Method Not Allowed")
	StatusAccessDenied     = NewRPCStatus(30022, "Access Denied")
)
