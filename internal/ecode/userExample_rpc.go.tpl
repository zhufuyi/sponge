package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// {{.TableNameCamelFCL}} business-level rpc error codes.
// the _{{.TableNameCamelFCL}}NO value range is 1~100, if the same error code is used, it will cause panic.
var (
	_{{.TableNameCamelFCL}}NO       = 2
	_{{.TableNameCamelFCL}}Name     = "{{.TableNameCamelFCL}}"
	_{{.TableNameCamelFCL}}BaseCode = errcode.RCode(_{{.TableNameCamelFCL}}NO)

	StatusCreate{{.TableNameCamel}}     = errcode.NewRPCStatus(_{{.TableNameCamelFCL}}BaseCode+1, "failed to create "+_{{.TableNameCamelFCL}}Name)
	StatusDeleteBy{{.ColumnNameCamel}}{{.TableNameCamel}} = errcode.NewRPCStatus(_{{.TableNameCamelFCL}}BaseCode+2, "failed to delete "+_{{.TableNameCamelFCL}}Name)
	StatusUpdateBy{{.ColumnNameCamel}}{{.TableNameCamel}} = errcode.NewRPCStatus(_{{.TableNameCamelFCL}}BaseCode+3, "failed to update "+_{{.TableNameCamelFCL}}Name)
	StatusGetBy{{.ColumnNameCamel}}{{.TableNameCamel}}    = errcode.NewRPCStatus(_{{.TableNameCamelFCL}}BaseCode+4, "failed to get "+_{{.TableNameCamelFCL}}Name+" details")
	StatusList{{.TableNameCamel}}       = errcode.NewRPCStatus(_{{.TableNameCamelFCL}}BaseCode+5, "failed to list of "+_{{.TableNameCamelFCL}}Name)

	// error codes are globally unique, adding 1 to the previous error code
)
