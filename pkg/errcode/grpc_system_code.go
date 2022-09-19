package errcode

// nolint
// rpc系统级别错误码，有status前缀
var (
	StatusSuccess = NewGRPCStatus(0, "ok")

	StatusInvalidParams       = NewGRPCStatus(300001, "参数错误")
	StatusUnauthorized        = NewGRPCStatus(300002, "认证错误")
	StatusInternalServerError = NewGRPCStatus(300003, "服务内部错误")
	StatusNotFound            = NewGRPCStatus(300004, "资源不存在")
	StatusAlreadyExists       = NewGRPCStatus(300005, "资源已存在")
	StatusTimeout             = NewGRPCStatus(300006, "超时")
	StatusTooManyRequests     = NewGRPCStatus(300007, "请求过多")
	StatusForbidden           = NewGRPCStatus(300008, "拒绝访问")
	StatusLimitExceed         = NewGRPCStatus(300009, "访问限制")

	StatusDeadlineExceeded = NewGRPCStatus(300010, "已超过最后期限")
	StatusAccessDenied     = NewGRPCStatus(300011, "拒绝访问")
	StatusMethodNotAllowed = NewGRPCStatus(300012, "不允许使用的方法")
)
