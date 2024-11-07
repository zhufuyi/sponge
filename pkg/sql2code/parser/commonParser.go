package parser

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/jinzhu/inflection"
)

// CrudInfo crud info for cache, dao, handler, service, protobuf, error
type CrudInfo struct {
	TableNameCamel          string `json:"tableNameCamel"`          // camel case, example: FooBar
	TableNameCamelFCL       string `json:"tableNameCamelFCL"`       // camel case and first character lower, example: fooBar
	TableNamePluralCamel    string `json:"tableNamePluralCamel"`    // plural, camel case, example: FooBars
	TableNamePluralCamelFCL string `json:"tableNamePluralCamelFCL"` // plural, camel case, example: fooBars

	ColumnName               string `json:"columnName"`               // column name, example: first_name
	ColumnNameCamel          string `json:"columnNameCamel"`          // column name, camel case, example: FirstName
	ColumnNameCamelFCL       string `json:"columnNameCamelFCL"`       // column name, camel case and first character lower, example: firstName
	ColumnNamePluralCamel    string `json:"columnNamePluralCamel"`    // column name, plural, camel case, example: FirstNames
	ColumnNamePluralCamelFCL string `json:"columnNamePluralCamelFCL"` // column name, plural, camel case and first character lower, example: firstNames

	GoType       string `json:"goType"`       // go type, example: string, uint64
	GoTypeFCU    string `json:"goTypeFCU"`    // go type, first character upper, example: String, Uint64
	ProtoType    string `json:"protoType"`    // proto type, example: string, uint64
	IsStringType bool   `json:"isStringType"` // go type is string or not

	PrimaryKeyColumnName string `json:"PrimaryKeyColumnName"` // primary key, example: id
	IsCommonType         bool   `json:"isCommonType"`         // custom primary key name and type
	IsStandardPrimaryKey bool   `json:"isStandardPrimaryKey"` // standard primary key id
}

func isDesiredGoType(t string) bool {
	switch t {
	case "string", "uint64", "int64", "uint", "int", "uint32", "int32": //nolint
		return true
	}
	return false
}

func setCrudInfo(field tmplField) *CrudInfo {
	primaryKeyName := ""
	if field.IsPrimaryKey {
		primaryKeyName = field.ColName
	}
	pluralName := inflection.Plural(field.Name)

	return &CrudInfo{
		ColumnName:               field.ColName,
		ColumnNameCamel:          field.Name,
		ColumnNameCamelFCL:       customFirstLetterToLower(field.Name),
		ColumnNamePluralCamel:    customEndOfLetterToLower(field.Name, pluralName),
		ColumnNamePluralCamelFCL: customFirstLetterToLower(customEndOfLetterToLower(field.Name, pluralName)),
		GoType:                   field.GoType,
		GoTypeFCU:                firstLetterToUpper(field.GoType),
		ProtoType:                simpleGoTypeToProtoType(field.GoType),
		IsStringType:             field.GoType == "string",
		PrimaryKeyColumnName:     primaryKeyName,
		IsStandardPrimaryKey:     field.ColName == "id",
	}
}

func newCrudInfo(data tmplData) *CrudInfo {
	if len(data.Fields) == 0 {
		return nil
	}

	var info *CrudInfo
	for _, field := range data.Fields {
		if field.IsPrimaryKey {
			info = setCrudInfo(field)
			break
		}
	}

	// if not found primary key, find the first xxx_id column as primary key
	if info == nil {
		for _, field := range data.Fields {
			if strings.HasSuffix(field.ColName, "_id") && isDesiredGoType(field.GoType) { // xxx_id
				info = setCrudInfo(field)
				break
			}
		}
	}

	// if not found xxx_id field, use the first field of integer or string type
	if info == nil {
		for _, field := range data.Fields {
			if isDesiredGoType(field.GoType) {
				info = setCrudInfo(field)
				break
			}
		}
	}

	// use the first column as primary key
	if info == nil {
		info = setCrudInfo(data.Fields[0])
	}

	info.TableNameCamel = data.TableName
	info.TableNameCamelFCL = data.TName

	pluralName := inflection.Plural(data.TableName)
	info.TableNamePluralCamel = customEndOfLetterToLower(data.TableName, pluralName)
	info.TableNamePluralCamelFCL = customFirstLetterToLower(customEndOfLetterToLower(data.TableName, pluralName))

	return info
}

func (info *CrudInfo) getCode() string {
	if info == nil {
		return ""
	}
	pkData, _ := json.Marshal(info)
	return string(pkData)
}

func (info *CrudInfo) CheckCommonType() bool {
	if info == nil {
		return false
	}
	return info.IsCommonType
}

func (info *CrudInfo) isIDPrimaryKey() bool {
	if info == nil {
		return false
	}
	if info.ColumnName == "id" && (info.GoType == "uint64" ||
		info.GoType == "int64" ||
		info.GoType == "uint" ||
		info.GoType == "int" ||
		info.GoType == "uint32" ||
		info.GoType == "int32") {
		return true
	}
	return false
}

func (info *CrudInfo) GetGRPCProtoValidation() string {
	if info == nil {
		return ""
	}
	if info.ProtoType == "string" {
		return `[(validate.rules).string.min_len = 1]`
	}
	return fmt.Sprintf(`[(validate.rules).%s.gt = 0]`, info.ProtoType)
}

func (info *CrudInfo) GetWebProtoValidation() string {
	if info == nil {
		return ""
	}
	if info.ProtoType == "string" {
		return fmt.Sprintf(`[(validate.rules).string.min_len = 1, (tagger.tags) = "uri:\"%s\""]`, info.ColumnNameCamelFCL)
	}
	return fmt.Sprintf(`[(validate.rules).%s.gt = 0, (tagger.tags) = "uri:\"%s\""]`, info.ProtoType, info.ColumnNameCamelFCL)
}

func getCommonHandlerStructCodes(data tmplData, jsonNamedType int) (string, error) {
	newFields := []tmplField{}
	for _, field := range data.Fields {
		if jsonNamedType == 0 { // snake case
			field.JSONName = customToSnake(field.ColName)
		} else {
			field.JSONName = customToCamel(field.ColName) // camel case (default)
		}
		newFields = append(newFields, field)
	}
	data.Fields = newFields

	postStructCode, err := tmplExecuteWithFilter(data, handlerCreateStructCommonTmpl)
	if err != nil {
		return "", fmt.Errorf("handlerCreateStructTmpl error: %v", err)
	}

	putStructCode, err := tmplExecuteWithFilter(data, handlerUpdateStructCommonTmpl, columnID)
	if err != nil {
		return "", fmt.Errorf("handlerUpdateStructTmpl error: %v", err)
	}

	getStructCode, err := tmplExecuteWithFilter(data, handlerDetailStructCommonTmpl, columnID, columnCreatedAt, columnUpdatedAt)
	if err != nil {
		return "", fmt.Errorf("handlerDetailStructTmpl error: %v", err)
	}

	return postStructCode + putStructCode + getStructCode, nil
}

func getCommonServiceStructCode(data tmplData) (string, error) {
	builder := strings.Builder{}
	err := serviceStructCommonTmpl.Execute(&builder, data)
	if err != nil {
		return "", err
	}
	code := builder.String()

	serviceCreateStructCode, err := tmplExecuteWithFilter(data, serviceCreateStructCommonTmpl)
	if err != nil {
		return "", fmt.Errorf("handle serviceCreateStructTmpl error: %v", err)
	}
	serviceCreateStructCode = strings.ReplaceAll(serviceCreateStructCode, "ID:", "Id:")

	serviceUpdateStructCode, err := tmplExecuteWithFilter(data, serviceUpdateStructCommonTmpl, columnID)
	if err != nil {
		return "", fmt.Errorf("handle serviceUpdateStructTmpl error: %v", err)
	}
	serviceUpdateStructCode = strings.ReplaceAll(serviceUpdateStructCode, "ID:", "Id:")

	code = strings.ReplaceAll(code, "// serviceCreateStructCode", serviceCreateStructCode)
	code = strings.ReplaceAll(code, "// serviceUpdateStructCode", serviceUpdateStructCode)

	return code, nil
}

func getCommonProtoFileCode(data tmplData, jsonNamedType int, isWebProto bool, isExtendedAPI bool) (string, error) {
	data.Fields = goTypeToProto(data.Fields, jsonNamedType, true)

	var err error
	builder := strings.Builder{}
	if isWebProto {
		if isExtendedAPI {
			err = protoFileForWebCommonTmpl.Execute(&builder, data)
		} else {
			err = protoFileForSimpleWebCommonTmpl.Execute(&builder, data)
		}
		if err != nil {
			return "", err
		}
	} else {
		if isExtendedAPI {
			err = protoFileCommonTmpl.Execute(&builder, data)
		} else {
			err = protoFileSimpleCommonTmpl.Execute(&builder, data)
		}
		if err != nil {
			return "", err
		}
	}
	code := builder.String()

	protoMessageCreateCode, err := tmplExecuteWithFilter2(data, protoMessageCreateCommonTmpl)
	if err != nil {
		return "", fmt.Errorf("handle protoMessageCreateCommonTmpl error: %v", err)
	}

	protoMessageUpdateCode, err := tmplExecuteWithFilter2(data, protoMessageUpdateCommonTmpl, columnID)
	if err != nil {
		return "", fmt.Errorf("handle protoMessageUpdateCommonTmpl error: %v", err)
	}
	if !isWebProto {
		srcStr := fmt.Sprintf(`, (tagger.tags) = "uri:\"%s\""`, getProtoFieldName(data.Fields))
		protoMessageUpdateCode = strings.ReplaceAll(protoMessageUpdateCode, srcStr, "")
	}

	protoMessageDetailCode, err := tmplExecuteWithFilter2(data, protoMessageDetailCommonTmpl, columnID, columnCreatedAt, columnUpdatedAt)
	if err != nil {
		return "", fmt.Errorf("handle protoMessageDetailCommonTmpl error: %v", err)
	}

	code = strings.ReplaceAll(code, "// protoMessageCreateCode", protoMessageCreateCode)
	code = strings.ReplaceAll(code, "// protoMessageUpdateCode", protoMessageUpdateCode)
	code = strings.ReplaceAll(code, "// protoMessageDetailCode", protoMessageDetailCode)
	code = strings.ReplaceAll(code, "*time.Time", "int64")
	code = strings.ReplaceAll(code, "time.Time", "int64")
	code = strings.ReplaceAll(code, "left_curly_bracket", "{")
	code = strings.ReplaceAll(code, "right_curly_bracket", "}")

	code = adaptedDbType2(data, isWebProto, code)

	return code, nil
}

func tmplExecuteWithFilter2(data tmplData, tmpl *template.Template, reservedColumns ...string) (string, error) {
	var newFields = []tmplField{}
	for _, field := range data.Fields {
		if isIgnoreFields(field.ColName, reservedColumns...) {
			continue
		}
		newFields = append(newFields, field)
	}
	data.Fields = newFields

	builder := strings.Builder{}
	err := tmpl.Execute(&builder, data)
	if err != nil {
		return "", fmt.Errorf("tmpl.Execute error: %v", err)
	}
	return builder.String(), nil
}

// nolint
func simpleGoTypeToProtoType(goType string) string {
	var protoType string
	switch goType {
	case "int", "int32":
		protoType = "int32"
	case "uint", "uint32":
		protoType = "uint32"
	case "int64":
		protoType = "int64"
	case "uint64":
		protoType = "uint64"
	case "string":
		protoType = "string"
	case "time.Time", "*time.Time":
		protoType = "string"
	case "float32":
		protoType = "float"
	case "float64":
		protoType = "double"
	case goTypeInts, "[]int64":
		protoType = "repeated int64"
	case "[]int32":
		protoType = "repeated int32"
	case "[]byte":
		protoType = "string"
	case goTypeStrings:
		protoType = "repeated string"
	case jsonTypeName:
		protoType = "string"
	default:
		protoType = "string"
	}
	return protoType
}

func adaptedDbType2(data tmplData, isWebProto bool, code string) string {
	if isWebProto {
		code = replaceProtoMessageFieldCode(code, webDefaultProtoMessageFieldCodes)
	} else {
		code = replaceProtoMessageFieldCode(code, grpcDefaultProtoMessageFieldCodes)
	}

	if data.ProtoSubStructs != "" {
		code += "\n" + data.ProtoSubStructs
	}

	return code
}

func firstLetterToUpper(str string) string {
	if len(str) == 0 {
		return str
	}

	if (str[0] >= 'A' && str[0] <= 'Z') || (str[0] >= 'a' && str[0] <= 'z') {
		return strings.ToUpper(str[:1]) + str[1:]
	}

	return str
}

func customFirstLetterToLower(str string) string {
	str = firstLetterToLower(str)

	if len(str) == 2 {
		if str == "iD" {
			str = "id"
		} else if str == "iP" {
			str = "ip"
		}
	} else if len(str) == 3 {
		if str == "iDs" {
			str = "ids"
		} else if str == "iPs" {
			str = "ips"
		}
	}

	return str
}

func customEndOfLetterToLower(srcStr string, str string) string {
	l := len(str) - len(srcStr)
	if l == 1 {
		if str[len(str)-1] == 'S' {
			return str[:len(str)-1] + "s"
		}
	} else if l == 2 {
		if str[len(str)-2:] == "ES" {
			return str[:len(str)-2] + "es"
		}
	}

	return str
}
