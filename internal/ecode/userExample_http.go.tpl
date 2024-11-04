package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// {{.TableNameCamelFCL}} business-level http error codes.
// the {{.TableNameCamelFCL}}NO value range is 1~100, if the same error code is used, it will cause panic.
var (
	{{.TableNameCamelFCL}}NO       = 1
	{{.TableNameCamelFCL}}Name     = "{{.TableNameCamelFCL}}"
	{{.TableNameCamelFCL}}BaseCode = errcode.HCode({{.TableNameCamelFCL}}NO)

	ErrCreate{{.TableNameCamel}}     = errcode.NewError({{.TableNameCamelFCL}}BaseCode+1, "failed to create "+{{.TableNameCamelFCL}}Name)
	ErrDeleteBy{{.ColumnNameCamel}}{{.TableNameCamel}} = errcode.NewError({{.TableNameCamelFCL}}BaseCode+2, "failed to delete "+{{.TableNameCamelFCL}}Name)
	ErrUpdateBy{{.ColumnNameCamel}}{{.TableNameCamel}} = errcode.NewError({{.TableNameCamelFCL}}BaseCode+3, "failed to update "+{{.TableNameCamelFCL}}Name)
	ErrGetBy{{.ColumnNameCamel}}{{.TableNameCamel}}    = errcode.NewError({{.TableNameCamelFCL}}BaseCode+4, "failed to get "+{{.TableNameCamelFCL}}Name+" details")
	ErrList{{.TableNameCamel}}       = errcode.NewError({{.TableNameCamelFCL}}BaseCode+5, "failed to list of "+{{.TableNameCamelFCL}}Name)

	// error codes are globally unique, adding 1 to the previous error code
)
