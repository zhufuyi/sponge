package errcode

// nolint
// 服务级别错误码，有Err前缀
var (
// ErrUserCreate = NewError(200101, "创建用户失败")
// ErrUserDelete = NewError(200102, "删除用户失败")
// ErrUserUpdate = NewError(200103, "更新用户失败")
// ErrUserGet    = NewError(200104, "获取用户失败")
)

// HCode 根据编号生成200000~300000之间的错误码
func HCode(NO int) int {
	if NO > 1000 {
		panic("NO must be < 1000")
	}
	return 200000 + NO*100
}
