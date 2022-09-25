package errcode

// nolint
// 服务级别错误码，有Err前缀
var (
// ErrUserCreate = NewError(HCode(1)+1, "创建用户失败")	// 200101
// ErrUserDelete = NewError(HCode(1)+2, "删除用户失败")	// 200102
// ErrUserUpdate = NewError(HCode(1)+3, "更新用户失败")	// 200103
// ErrUserGet    = NewError(HCode(1)+4, "获取用户失败") 		// 200104
)

// HCode 根据编号生成200000~300000之间的错误码
func HCode(NO int) int {
	if NO > 1000 {
		panic("NO must be < 1000")
	}
	return 200000 + NO*100
}
