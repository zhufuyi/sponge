package errcode

// nolint
// rpc系统级别错误码，有status前缀
var (
	StatusSuccess = NewRPCStatus(0, "ok")

	StatusInvalidParams       = NewRPCStatus(300001, "参数错误")
	StatusUnauthorized        = NewRPCStatus(300002, "认证错误")
	StatusInternalServerError = NewRPCStatus(300003, "服务内部错误")
	StatusNotFound            = NewRPCStatus(300004, "资源不存在")
	StatusAlreadyExists       = NewRPCStatus(300005, "资源已存在")
	StatusTimeout             = NewRPCStatus(300006, "访问超时")
	StatusTooManyRequests     = NewRPCStatus(300007, "请求过多")
	StatusForbidden           = NewRPCStatus(300008, "拒绝访问")
	StatusLimitExceed         = NewRPCStatus(300009, "访问限制")

	StatusDeadlineExceeded   = NewRPCStatus(300010, "已超过最后期限")
	StatusAccessDenied       = NewRPCStatus(300011, "拒绝访问")
	StatusMethodNotAllowed   = NewRPCStatus(300012, "不允许使用的方法")
	StatusServiceUnavailable = NewRPCStatus(300013, "服务不可用")
)
