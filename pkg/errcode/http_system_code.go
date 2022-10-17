package errcode

// nolint
// http系统级别错误码，无Err前缀
var (
	Success             = NewError(0, "ok")
	InvalidParams       = NewError(100001, "参数错误")
	Unauthorized        = NewError(100002, "认证错误")
	InternalServerError = NewError(100003, "服务内部错误")
	NotFound            = NewError(100004, "资源不存在")
	AlreadyExists       = NewError(100005, "资源已存在")
	Timeout             = NewError(100006, "超时")
	TooManyRequests     = NewError(100007, "请求过多")
	Forbidden           = NewError(100008, "拒绝访问")
	LimitExceed         = NewError(100009, "访问限制")

	DeadlineExceeded         = NewError(100010, "已超过最后期限")
	AccessDenied             = NewError(100011, "拒绝访问")
	MethodNotAllowed         = NewError(100012, "不允许使用的方法")
	MethodServiceUnavailable = NewError(100013, "服务不可用")
)
