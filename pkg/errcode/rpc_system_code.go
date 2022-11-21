package errcode

// nolint
// rpc system level error code with status prefix, error code range 30000~40000
var (
	StatusSuccess = NewRPCStatus(0, "ok")

	StatusInvalidParams       = NewRPCStatus(30001, "Invalid Parameter")
	StatusUnauthorized        = NewRPCStatus(30002, "Unauthorized")
	StatusInternalServerError = NewRPCStatus(30003, "Internal Server Error")
	StatusNotFound            = NewRPCStatus(30004, "Not Found")
	StatusAlreadyExists       = NewRPCStatus(30005, "Conflict")
	StatusTimeout             = NewRPCStatus(30006, "Request Timeout")
	StatusTooManyRequests     = NewRPCStatus(30007, "Too Many Requests")
	StatusForbidden           = NewRPCStatus(30008, "Forbidden")
	StatusLimitExceed         = NewRPCStatus(30009, "Limit Exceed")
	StatusDeadlineExceeded    = NewRPCStatus(30010, "Deadline Exceeded")
	StatusAccessDenied        = NewRPCStatus(30011, "Access Denied")
	StatusMethodNotAllowed    = NewRPCStatus(30012, "Method Not Allowed")
	StatusServiceUnavailable  = NewRPCStatus(30013, "Service Unavailable")
)
