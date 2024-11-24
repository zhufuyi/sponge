package parser

import (
	"encoding/json"
	"fmt"

	"github.com/huandu/xstrings"
	"github.com/jinzhu/inflection"
)

// TableInfo is the struct for extend template
type TableInfo struct {
	TableNamePrefix string // table name prefix, example: t_

	TableName               string // original table name, example: foo_bar
	TableNameCamel          string // camel case, example: FooBar
	TableNameCamelFCL       string // camel case and first character lower, example: fooBar
	TableNamePluralCamel    string // plural, camel case, example: FooBars
	TableNamePluralCamelFCL string // plural, camel case and first character lower, example: fooBars
	TableNameSnake          string // snake case, example: foo_bar

	TableComment string // table comment

	Columns    []Field     // columns of the table
	PrimaryKey *PrimaryKey // primary key information

	DBDriver string // database driver, example: mysql, postgresql, sqlite3, mongodb

	ColumnSubStructure string // column sub structure for model
	ColumnSubMessage   string // sub message for protobuf
}

// Field is the struct for column information
type Field struct {
	ColumnName         string // original column name, example: foo_bar
	ColumnNameCamel    string // first character lower, example: FooBar
	ColumnNameCamelFCL string // first character lower, example: fooBar

	ColumnComment string // column comment
	IsPrimaryKey  bool   // is primary key

	GoType string // convert to go type
	Tag    string // tag for model struct field, default gorm tag
}

// PrimaryKey is the struct for primary key information, it used for generate CRUD code
type PrimaryKey struct {
	Name               string // primary key name, example: foo_bar
	NameCamel          string // primary key name, camel case, example: FooBar
	NameCamelFCL       string // primary key name, camel case and first character lower, example: fooBar
	NamePluralCamel    string // primary key name, plural, camel case, example: FooBars
	NamePluralCamelFCL string // primary key name, plural, camel case and first character lower, example: fooBars

	GoType    string // go type, example:  int, string
	GoTypeFCU string // go type, first character upper, example: Int64, String

	IsStringType bool // go type is string or not
}

func newTableInfo(data tmplData) TableInfo {
	pluralName := inflection.Plural(data.TableName)
	return TableInfo{
		TableNamePrefix:         data.TableNamePrefix,
		TableName:               data.RawTableName,
		TableNameCamel:          data.TableName,
		TableNameCamelFCL:       data.TName,
		TableNamePluralCamel:    customEndOfLetterToLower(data.TableName, pluralName),
		TableNamePluralCamelFCL: customFirstLetterToLower(customEndOfLetterToLower(data.TableName, pluralName)),
		TableNameSnake:          xstrings.ToSnakeCase(data.TName),
		TableComment:            data.Comment,
		Columns:                 getColumns(data.Fields),
		PrimaryKey:              getPrimaryKeyInfo(data.CrudInfo),
		DBDriver:                data.DBDriver,
		ColumnSubStructure:      data.SubStructs,
		ColumnSubMessage:        data.ProtoSubStructs,
	}
}

func (table TableInfo) getCode() []byte {
	code, err := json.Marshal(&table)
	if err != nil {
		fmt.Printf("table: %v, json.Marshal error: %v\n", table.TableName, err)
	}
	return code
}

func getColumns(fields []tmplField) []Field {
	var columns []Field

	for _, field := range fields {
		columns = append(columns, Field{
			ColumnName:         field.ColName,
			ColumnNameCamel:    field.Name,
			ColumnNameCamelFCL: customFirstLetterToLower(field.Name),
			ColumnComment:      field.Comment,
			IsPrimaryKey:       field.IsPrimaryKey,
			GoType:             field.GoType,
			Tag:                field.Tag,
		})
	}

	return columns
}

func getPrimaryKeyInfo(info *CrudInfo) *PrimaryKey {
	if info == nil {
		return nil
	}
	return &PrimaryKey{
		Name:               info.ColumnName,
		NameCamel:          info.ColumnNameCamel,
		NameCamelFCL:       info.ColumnNameCamelFCL,
		NamePluralCamel:    info.ColumnNamePluralCamel,
		NamePluralCamelFCL: info.ColumnNamePluralCamelFCL,
		GoType:             info.GoType,
		GoTypeFCU:          info.GoTypeFCU,
		IsStringType:       info.IsStringType,
	}
}

// UnMarshalTableInfo unmarshal the json data to TableInfo struct
func UnMarshalTableInfo(data string) (map[string]interface{}, error) {
	info := map[string]interface{}{}
	err := json.Unmarshal([]byte(data), &info)
	if err != nil {
		return info, err
	}
	return info, nil
}
