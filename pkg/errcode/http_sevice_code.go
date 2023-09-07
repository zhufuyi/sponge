package errcode

// HCode Generate an error code between 20000 and 30000 according to the number
//
// http service level error code, Err prefix, example.
//
// var (
// ErrUserCreate = NewError(HCode(1)+1, "failed to create user")		// 20101
// ErrUserDelete = NewError(HCode(1)+2, "failed to delete user")		// 20102
// ErrUserUpdate = NewError(HCode(1)+3, "failed to update user")		// 20103
// ErrUserGet    = NewError(HCode(1)+4, "failed to get user details")	// 20104
// )
func HCode(num int) int {
	if num > 99 || num < 1 {
		panic("num range must be between 0 to 100")
	}
	return 20000 + num*100
}
