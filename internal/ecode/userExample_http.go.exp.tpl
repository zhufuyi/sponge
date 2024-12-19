package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// {{.TableNameCamelFCL}} business-level http error codes.
// the {{.TableNameCamelFCL}}NO value range is 1~100, if the same error code is used, it will cause panic.
var (
	{{.TableNameCamelFCL}}NO       = 78
	{{.TableNameCamelFCL}}Name     = "{{.TableNameCamelFCL}}"
	{{.TableNameCamelFCL}}BaseCode = errcode.HCode({{.TableNameCamelFCL}}NO)

	ErrCreate{{.TableNameCamel}}     = errcode.NewError({{.TableNameCamelFCL}}BaseCode+1, "failed to create "+{{.TableNameCamelFCL}}Name)
	ErrDeleteBy{{.ColumnNameCamel}}{{.TableNameCamel}} = errcode.NewError({{.TableNameCamelFCL}}BaseCode+2, "failed to delete "+{{.TableNameCamelFCL}}Name)
	ErrUpdateBy{{.ColumnNameCamel}}{{.TableNameCamel}} = errcode.NewError({{.TableNameCamelFCL}}BaseCode+3, "failed to update "+{{.TableNameCamelFCL}}Name)
	ErrGetBy{{.ColumnNameCamel}}{{.TableNameCamel}}    = errcode.NewError({{.TableNameCamelFCL}}BaseCode+4, "failed to get "+{{.TableNameCamelFCL}}Name+" details")
	ErrList{{.TableNameCamel}}       = errcode.NewError({{.TableNameCamelFCL}}BaseCode+5, "failed to list of "+{{.TableNameCamelFCL}}Name)

	ErrDeleteBy{{.ColumnNamePluralCamel}}{{.TableNameCamel}}    = errcode.NewError({{.TableNameCamelFCL}}BaseCode+6, "failed to delete by batch {{.ColumnNamePluralCamelFCL}} "+{{.TableNameCamelFCL}}Name)
	ErrGetByCondition{{.TableNameCamel}} = errcode.NewError({{.TableNameCamelFCL}}BaseCode+7, "failed to get "+{{.TableNameCamelFCL}}Name+" details by conditions")
	ErrListBy{{.ColumnNamePluralCamel}}{{.TableNameCamel}}      = errcode.NewError({{.TableNameCamelFCL}}BaseCode+8, "failed to list by batch {{.ColumnNamePluralCamelFCL}} "+{{.TableNameCamelFCL}}Name)
	ErrListByLast{{.ColumnNameCamel}}{{.TableNameCamel}}   = errcode.NewError({{.TableNameCamelFCL}}BaseCode+9, "failed to list by last {{.ColumnNameCamelFCL}} "+{{.TableNameCamelFCL}}Name)

	// error codes are globally unique, adding 1 to the previous error code
)
