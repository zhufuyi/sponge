package errcode

import "google.golang.org/grpc/codes"

// nolint
// rpc服务级别错误码，有Status前缀
var (
// StatusUserCreate = NewGRPCStatus(GCode(1)+1, "创建用户失败")		// 400101
// StatusUserDelete = NewGRPCStatus(GCode(1)+2, "删除用户失败")		// 400102
// StatusUserUpdate = NewGRPCStatus(GCode(1)+3, "更新用户失败")	// 400103
// StatusUserGet    = NewGRPCStatus(GCode(1)+4, "获取用户失败")		// 400104
)

// GCode 根据编号生成400000~500000之间的错误码
func GCode(NO int) codes.Code {
	if NO > 1000 {
		panic("NO must be < 1000")
	}
	return codes.Code(400000 + NO*100)
}
