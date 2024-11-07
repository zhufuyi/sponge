package errcode

// http system level error code, error code range 10000~20000
var (
	Success             = NewError(0, "ok")
	InvalidParams       = NewError(100001, "Invalid Parameter")
	Unauthorized        = NewError(100002, "Unauthorized")
	InternalServerError = NewError(100003, "Internal Server Error")
	NotFound            = NewError(100004, "Not Found")
	Timeout             = NewError(100006, "Request Timeout")
	TooManyRequests     = NewError(100007, "Too Many Requests")
	Forbidden           = NewError(100008, "Forbidden")
	LimitExceed         = NewError(100009, "Limit Exceed")
	DeadlineExceeded    = NewError(100010, "Deadline Exceeded")
	AccessDenied        = NewError(100011, "Access Denied")
	MethodNotAllowed    = NewError(100012, "Method Not Allowed")
	ServiceUnavailable  = NewError(100013, "Service Unavailable")

	Canceled           = NewError(100014, "Canceled")
	Unknown            = NewError(100015, "Unknown")
	PermissionDenied   = NewError(100016, "Permission Denied")
	ResourceExhausted  = NewError(100017, "Resource Exhausted")
	FailedPrecondition = NewError(100018, "Failed Precondition")
	Aborted            = NewError(100019, "Aborted")
	OutOfRange         = NewError(100020, "Out Of Range")
	Unimplemented      = NewError(100021, "Unimplemented")
	DataLoss           = NewError(100022, "Data Loss")

	StatusBadGateway = NewError(100023, "Bad Gateway")

	// Deprecated: use Conflict instead
	AlreadyExists = NewError(100005, "Already Exists")
	Conflict      = NewError(100409, "Conflict")
	TooEarly      = NewError(100425, "Too Early")
)
