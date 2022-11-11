package errcode

import "google.golang.org/grpc/codes"

// RCode Generate an error code between 40000 and 50000 according to the number
//
// rpc service level error code, status prefix, example.
//
//	var (
//		StatusUserCreate = NewRPCStatus(RCode(1)+1, "failed to create user")		// 40101
//		StatusUserDelete = NewRPCStatus(RCode(1)+2, "failed to delete user")		// 40102
//		StatusUserUpdate = NewRPCStatus(RCode(1)+3, "failed to update user")		// 40103
//		StatusUserGet    = NewRPCStatus(RCode(1)+4, "failed to get user details")	// 40104
//	)
func RCode(NO int) codes.Code {
	if NO > 99 {
		panic("NO must be < 100")
	}
	return codes.Code(40000 + NO*100)
}
