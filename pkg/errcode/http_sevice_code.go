package errcode

// HCode Generate an error code between 20000 and 30000 according to the number
//
// http service level error code, Err prefix, example.
//
// var (
// ErrUserCreate = NewError(HCode(1)+1, "failed to create user")		// 200101
// ErrUserDelete = NewError(HCode(1)+2, "failed to delete user")			// 200102
// ErrUserUpdate = NewError(HCode(1)+3, "failed to update user")		// 200103
// ErrUserGet    = NewError(HCode(1)+4, "failed to get user details")	// 200104
// )
func HCode(NO int) int {
	if NO > 99 || NO < 1 {
		panic("NO range must be between 0 to 100")
	}
	return 200000 + NO*100
}
