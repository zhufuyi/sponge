package errcode

import "google.golang.org/grpc/codes"

// nolint
// rpc服务级别错误码，有Status前缀
var (
// StatusUserCreate = NewGRPCStatus(400101, "创建用户失败")
// StatusUserDelete = NewGRPCStatus(400102, "删除用户失败")
// StatusUserUpdate = NewGRPCStatus(400103, "更新用户失败")
// StatusUserGet    = NewGRPCStatus(400104, "获取用户失败")
)

// GCode 根据编号生成400000~500000之间的错误码
func GCode(NO int) codes.Code {
	if NO > 1000 {
		panic("NO must be < 1000")
	}
	return codes.Code(400000 + NO*100)
}
