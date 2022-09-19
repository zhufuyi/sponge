package errcode

// 系统级别错误码，无Err前缀
var (
	Success             = NewError(0, "ok")
	InvalidParams       = NewError(10001, "参数错误")
	Unauthorized        = NewError(10002, "认证错误")
	InternalServerError = NewError(10003, "服务内部错误")
	NotFound            = NewError(10004, "资源不存在")
	AlreadyExists       = NewError(10005, "资源已存在")
	Timeout             = NewError(10006, "超时")
	TooManyRequests     = NewError(10007, "请求过多")
	Forbidden           = NewError(10008, "拒绝访问")
	LimitExceed         = NewError(10009, "访问限制")
)
