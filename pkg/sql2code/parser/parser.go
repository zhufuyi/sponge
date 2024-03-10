// Package parser is a library that parses to go structures based on sql
// and generates the code needed based on the template.
package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"sort"
	"strings"
	"text/template"

	"github.com/blastrain/vitess-sqlparser/tidbparser/ast"
	"github.com/blastrain/vitess-sqlparser/tidbparser/dependency/mysql"
	"github.com/blastrain/vitess-sqlparser/tidbparser/dependency/types"
	"github.com/blastrain/vitess-sqlparser/tidbparser/parser"
	"github.com/huandu/xstrings"
	"github.com/jinzhu/inflection"
)

const (
	// TableName table name
	TableName = "__table_name__"
	// CodeTypeModel model code
	CodeTypeModel = "model"
	// CodeTypeJSON json code
	CodeTypeJSON = "json"
	// CodeTypeDAO update fields code
	CodeTypeDAO = "dao"
	// CodeTypeHandler handler request and respond code
	CodeTypeHandler = "handler"
	// CodeTypeProto proto file code
	CodeTypeProto = "proto"
	// CodeTypeService grpc service code
	CodeTypeService = "service"

	// DBDriverMysql mysql driver
	DBDriverMysql = "mysql"
	// DBDriverPostgresql postgresql driver
	DBDriverPostgresql = "postgresql"
	// DBDriverTidb tidb driver
	DBDriverTidb = "tidb"
	// DBDriverSqlite sqlite driver
	DBDriverSqlite = "sqlite"
	// DBDriverMongodb mongodb driver
	DBDriverMongodb = "mongodb"
)

// Codes content
type Codes struct {
	Model         []string // model code
	UpdateFields  []string // update fields code
	ModelJSON     []string // model json code
	HandlerStruct []string // handler request and respond code
}

// modelCodes model code
type modelCodes struct {
	Package    string
	ImportPath []string
	StructCode []string
}

// ParseSQL generate different usage codes based on sql
func ParseSQL(sql string, options ...Option) (map[string]string, error) {
	initTemplate()
	opt := parseOption(options)

	stmts, err := parser.New().Parse(sql, opt.Charset, opt.Collation)
	if err != nil {
		return nil, err
	}
	modelStructCodes := make([]string, 0, len(stmts))
	updateFieldsCodes := make([]string, 0, len(stmts))
	handlerStructCodes := make([]string, 0, len(stmts))
	protoFileCodes := make([]string, 0, len(stmts))
	serviceStructCodes := make([]string, 0, len(stmts))
	modelJSONCodes := make([]string, 0, len(stmts))
	importPath := make(map[string]struct{})
	tableNames := make([]string, 0, len(stmts))
	for _, stmt := range stmts {
		if ct, ok := stmt.(*ast.CreateTableStmt); ok {
			code, err2 := makeCode(ct, opt)
			if err2 != nil {
				return nil, err2
			}
			modelStructCodes = append(modelStructCodes, code.modelStruct)
			updateFieldsCodes = append(updateFieldsCodes, code.updateFields)
			handlerStructCodes = append(handlerStructCodes, code.handlerStruct)
			protoFileCodes = append(protoFileCodes, code.protoFile)
			serviceStructCodes = append(serviceStructCodes, code.serviceStruct)
			modelJSONCodes = append(modelJSONCodes, code.modelJSON)
			tableNames = append(tableNames, toCamel(ct.Table.Name.String()))
			for _, s := range code.importPaths {
				importPath[s] = struct{}{}
			}
		}
	}

	importPathArr := make([]string, 0, len(importPath))
	for s := range importPath {
		importPathArr = append(importPathArr, s)
	}
	sort.Strings(importPathArr)

	mc := modelCodes{
		Package:    opt.Package,
		ImportPath: importPathArr,
		StructCode: modelStructCodes,
	}
	modelCode, err := getModelCode(mc)
	if err != nil {
		return nil, err
	}

	var codesMap = map[string]string{
		CodeTypeModel:   modelCode,
		CodeTypeJSON:    strings.Join(modelJSONCodes, "\n\n"),
		CodeTypeDAO:     strings.Join(updateFieldsCodes, "\n\n"),
		CodeTypeHandler: strings.Join(handlerStructCodes, "\n\n"),
		CodeTypeProto:   strings.Join(protoFileCodes, "\n\n"),
		CodeTypeService: strings.Join(serviceStructCodes, "\n\n"),
		TableName:       strings.Join(tableNames, ", "),
	}

	return codesMap, nil
}

type tmplData struct {
	TableName       string
	TName           string
	NameFunc        bool
	RawTableName    string
	Fields          []tmplField
	Comment         string
	SubStructs      string // sub structs for model
	ProtoSubStructs string // sub structs for protobuf
	DBDriver        string
}

type tmplField struct {
	Name     string
	ColName  string
	GoType   string
	Tag      string
	Comment  string
	JSONName string
	DBDriver string
}

// ConditionZero type of condition 0, used in dao template code
func (t tmplField) ConditionZero() string {
	switch t.GoType {
	case "int8", "int16", "int32", "int64", "int", "uint8", "uint16", "uint32", "uint64", "uint", "float64", "float32", //nolint
		"sql.NullInt32", "sql.NullInt64", "sql.NullFloat64": //nolint
		return `!= 0`
	case "string", "sql.NullString": //nolint
		return `!= ""`
	case "time.Time", "*time.Time", "sql.NullTime": //nolint
		return `.IsZero() == false`
	case "[]byte", "[]string", "[]int", "interface{}": //nolint
		return `!= nil` //nolint
	case "bool": //nolint
		return `!= false /*Warning: if the value itself is false, can't be updated*/`
	}

	if t.DBDriver == DBDriverMongodb {
		if t.GoType == goTypeOID {
			return `!= primitive.NilObjectID`
		}
		if t.GoType == "*"+t.Name {
			return `!= nil`
		}
		if strings.Contains(t.GoType, "[]") {
			return `!= nil`
		}
	}

	return `!= ` + t.GoType
}

// GoZero type of 0, used in model to json template code
func (t tmplField) GoZero() string {
	switch t.GoType {
	case "int8", "int16", "int32", "int64", "int", "uint8", "uint16", "uint32", "uint64", "uint", "float64", "float32",
		"sql.NullInt32", "sql.NullInt64", "sql.NullFloat64":
		return `= 0`
	case "string", "sql.NullString":
		return `= "string"`
	case "time.Time", "*time.Time", "sql.NullTime":
		return `= "0000-01-00T00:00:00.000+08:00"`
	case "[]byte", "[]string", "[]int", "interface{}": //nolint
		return `= nil` //nolint
	case "bool": //nolint
		return `= false`
	}

	if t.DBDriver == DBDriverMongodb {
		if t.GoType == goTypeOID {
			return `= primitive.NilObjectID`
		}
		if t.GoType == "*"+t.Name {
			return `= nil`
		}
		if strings.Contains(t.GoType, "[]") {
			return `= nil`
		}
	}

	return `= ` + t.GoType
}

// GoTypeZero type of 0, used in service template code
func (t tmplField) GoTypeZero() string {
	switch t.GoType {
	case "int8", "int16", "int32", "int64", "int", "uint8", "uint16", "uint32", "uint64", "uint", "float64", "float32",
		"sql.NullInt32", "sql.NullInt64", "sql.NullFloat64":
		return `0`
	case "string", "sql.NullString":
		return `""`
	case "time.Time", "*time.Time", "sql.NullTime":
		return `0 /*time.Now().Second()*/`
	case "[]byte", "[]string", "[]int", "interface{}": //nolint
		return `nil` //nolint
	case "bool": //nolint
		return `false`
	}

	if t.DBDriver == DBDriverMongodb {
		if t.GoType == goTypeOID {
			return `primitive.NilObjectID`
		}
		if t.GoType == "*"+t.Name {
			return `nil` //nolint
		}
		if strings.Contains(t.GoType, "[]") {
			return `nil` //nolint
		}
	}

	return t.GoType
}

// AddOne counter
func (t tmplField) AddOne(i int) int {
	return i + 1
}

// AddOneWithTag counter and add id tag
func (t tmplField) AddOneWithTag(i int) string {
	if t.ColName == "id" {
		return fmt.Sprintf(`%d [(tagger.tags) = "uri:\"id\"" ]`, i+1)
	}

	return fmt.Sprintf("%d", i+1)
}

const (
	__mysqlModel__ = "__mysqlModel__" //nolint
	__type__       = "__type__"       //nolint
)

var replaceFields = map[string]string{
	__mysqlModel__: "ggorm.Model",
	__type__:       "",
}

const (
	columnID         = "id"
	_columnID        = "_id"
	columnCreatedAt  = "created_at"
	columnUpdatedAt  = "updated_at"
	columnDeletedAt  = "deleted_at"
	columnMysqlModel = __mysqlModel__
)

var ignoreColumns = map[string]struct{}{
	columnID:         {},
	columnCreatedAt:  {},
	columnUpdatedAt:  {},
	columnDeletedAt:  {},
	columnMysqlModel: {},
}

func isIgnoreFields(colName string, falseColumn ...string) bool {
	for _, v := range falseColumn {
		if colName == v {
			return false
		}
	}

	_, ok := ignoreColumns[colName]
	return ok
}

type codeText struct {
	importPaths   []string
	modelStruct   string
	modelJSON     string
	updateFields  string
	handlerStruct string
	protoFile     string
	serviceStruct string
}

// nolint
func makeCode(stmt *ast.CreateTableStmt, opt options) (*codeText, error) {
	importPath := make([]string, 0, 1)
	data := tmplData{
		TableName:    stmt.Table.Name.String(),
		RawTableName: stmt.Table.Name.String(),
		Fields:       make([]tmplField, 0, 1),
	}
	tablePrefix := opt.TablePrefix
	if tablePrefix != "" && strings.HasPrefix(data.TableName, tablePrefix) {
		data.NameFunc = true
		data.TableName = data.TableName[len(tablePrefix):]
	}
	if opt.ForceTableName || data.RawTableName != inflection.Plural(data.RawTableName) {
		data.NameFunc = true
	}
	switch opt.DBDriver {
	case DBDriverMongodb:
		if opt.JSONNamedType != 0 {
			SetJSONTagCamelCase()
		} else {
			SetJSONTagSnakeCase()
		}
	}

	data.TableName = toCamel(data.TableName)
	data.TName = firstLetterToLow(data.TableName)

	// find table comment
	for _, o := range stmt.Options {
		if o.Tp == ast.TableOptionComment {
			data.Comment = o.StrValue
			break
		}
	}

	isPrimaryKey := make(map[string]bool)
	for _, con := range stmt.Constraints {
		if con.Tp == ast.ConstraintPrimaryKey {
			isPrimaryKey[con.Keys[0].Column.String()] = true
		}
	}

	columnPrefix := opt.ColumnPrefix
	for _, col := range stmt.Cols {
		colName := col.Name.Name.String()
		goFieldName := colName
		if columnPrefix != "" && strings.HasPrefix(goFieldName, columnPrefix) {
			goFieldName = goFieldName[len(columnPrefix):]
		}
		jsonName := colName
		if opt.JSONNamedType != 0 {
			jsonName = xstrings.FirstRuneToLower(xstrings.ToCamelCase(colName)) // name type use camel case
		}
		field := tmplField{
			Name:     toCamel(goFieldName),
			ColName:  colName,
			JSONName: jsonName,
		}

		tags := make([]string, 0, 4)
		// make GORM's tag
		gormTag := strings.Builder{}
		gormTag.WriteString("column:")
		gormTag.WriteString(colName)
		if opt.GormType {
			gormTag.WriteString(";type:")
			switch opt.DBDriver {
			case DBDriverMysql, DBDriverTidb, DBDriverSqlite:
				gormTag.WriteString(col.Tp.InfoSchemaStr())
			case DBDriverPostgresql:
				gormTag.WriteString(opt.FieldTypes[colName])
			}
		}
		if isPrimaryKey[colName] {
			gormTag.WriteString(";primary_key")
		}
		isNotNull := false
		canNull := false
		for _, o := range col.Options {
			switch o.Tp {
			case ast.ColumnOptionPrimaryKey:
				if !isPrimaryKey[colName] {
					gormTag.WriteString(";primary_key")
					isPrimaryKey[colName] = true
				}
			case ast.ColumnOptionNotNull:
				isNotNull = true
			case ast.ColumnOptionAutoIncrement:
				gormTag.WriteString(";AUTO_INCREMENT")
			case ast.ColumnOptionDefaultValue:
				if value := getDefaultValue(o.Expr); value != "" {
					gormTag.WriteString(";default:")
					gormTag.WriteString(value)
				}
			case ast.ColumnOptionUniqKey:
				gormTag.WriteString(";unique")
			case ast.ColumnOptionNull:
				//gormTag.WriteString(";NULL")
				canNull = true
			case ast.ColumnOptionOnUpdate: // For Timestamp and Datetime only.
			case ast.ColumnOptionFulltext:
			case ast.ColumnOptionComment:
				field.Comment = o.Expr.GetDatum().GetString()
			default:
				//return "", nil, errors.Errorf(" unsupport option %d\n", o.Tp)
			}
		}

		field.DBDriver = opt.DBDriver
		switch opt.DBDriver {
		case DBDriverMongodb: // mongodb
			tags = append(tags, "bson", gormTag.String())
			if opt.JSONTag {
				if strings.ToLower(jsonName) == "_id" {
					jsonName = "id"
				}
				field.JSONName = jsonName
				tags = append(tags, "json", jsonName)
			}
			field.Tag = makeTagStr(tags)
			field.GoType = opt.FieldTypes[colName]
			if field.GoType == "time.Time" {
				importPath = append(importPath, "time")
			}

		default: // gorm
			if !isPrimaryKey[colName] && isNotNull {
				gormTag.WriteString(";NOT NULL")
			}
			tags = append(tags, "gorm", gormTag.String())

			if opt.JSONTag {
				tags = append(tags, "json", jsonName)
			}
			field.Tag = makeTagStr(tags)

			// get type in golang
			nullStyle := opt.NullStyle
			if !canNull {
				nullStyle = NullDisable
			}
			goType, pkg := mysqlToGoType(col.Tp, nullStyle)
			if pkg != "" {
				importPath = append(importPath, pkg)
			}
			field.GoType = goType
		}

		data.Fields = append(data.Fields, field)
	}
	if v, ok := opt.FieldTypes[SubStructKey]; ok {
		data.SubStructs = v
	}
	if v, ok := opt.FieldTypes[ProtoSubStructKey]; ok {
		data.ProtoSubStructs = v
	}
	data.DBDriver = opt.DBDriver

	updateFieldsCode, err := getUpdateFieldsCode(data, opt.IsEmbed)
	if err != nil {
		return nil, err
	}

	handlerStructCode, err := getHandlerStructCodes(data)
	if err != nil {
		return nil, err
	}

	modelStructCode, importPaths, err := getModelStructCode(data, importPath, opt.IsEmbed)
	if err != nil {
		return nil, err
	}

	modelJSONCode, err := getModelJSONCode(data)
	if err != nil {
		return nil, err
	}

	protoFileCode, err := getProtoFileCode(data, opt.IsWebProto)
	if err != nil {
		return nil, err
	}

	serviceStructCode, err := getServiceStructCode(data)
	if err != nil {
		return nil, err
	}

	return &codeText{
		importPaths:   importPaths,
		modelStruct:   modelStructCode,
		modelJSON:     modelJSONCode,
		updateFields:  updateFieldsCode,
		handlerStruct: handlerStructCode,
		protoFile:     protoFileCode,
		serviceStruct: serviceStructCode,
	}, nil
}

func getModelStructCode(data tmplData, importPaths []string, isEmbed bool) (string, []string, error) {
	// filter to ignore field fields
	var newFields = []tmplField{}
	var newImportPaths = []string{}
	if isEmbed {
		newFields = append(newFields, tmplField{
			Name:    __mysqlModel__,
			ColName: __mysqlModel__,
			GoType:  __type__,
			Tag:     `gorm:"embedded"`,
			Comment: "embed id and time\n",
		})

		isHaveTimeType := false
		for _, field := range data.Fields {
			if isIgnoreFields(field.ColName) {
				continue
			}
			newFields = append(newFields, field)
			if strings.Contains(field.GoType, "time.Time") {
				isHaveTimeType = true
			}
		}
		data.Fields = newFields

		// filter time package name
		if isHaveTimeType {
			newImportPaths = importPaths
		} else {
			for _, path := range importPaths {
				if path == "time" { //nolint
					continue
				}
				newImportPaths = append(newImportPaths, path)
			}
		}
		newImportPaths = append(newImportPaths, "github.com/zhufuyi/sponge/pkg/ggorm")
	} else {
		for i, field := range data.Fields {
			switch field.DBDriver {
			case DBDriverMongodb:
				if field.Name == "ID" {
					data.Fields[i].GoType = goTypeOID
					importPaths = append(importPaths, "go.mongodb.org/mongo-driver/bson/primitive")
				}

			default:
				if strings.Contains(field.GoType, "time.Time") {
					data.Fields[i].GoType = "*time.Time"
					continue
				}
				// force conversion of ID field to uint64 type
				if field.Name == "ID" {
					data.Fields[i].GoType = "uint64"
				}
			}
		}
		newImportPaths = importPaths
	}

	builder := strings.Builder{}
	err := modelStructTmpl.Execute(&builder, data)
	if err != nil {
		return "", nil, fmt.Errorf("modelStructTmpl.Execute error: %v", err)
	}
	code, err := format.Source([]byte(builder.String()))
	if err != nil {
		return "", nil, fmt.Errorf("modelStructTmpl format.Source error: %v", err)
	}
	structCode := string(code)
	// restore the real embedded fields
	if isEmbed {
		structCode = strings.ReplaceAll(structCode, __mysqlModel__, replaceFields[__mysqlModel__])
		structCode = strings.ReplaceAll(structCode, __type__, replaceFields[__type__])
	}

	if data.SubStructs != "" {
		structCode += data.SubStructs
	}
	if data.DBDriver == DBDriverMongodb {
		structCode = strings.ReplaceAll(structCode, `bson:"column:`, `bson:"`)
		structCode = strings.ReplaceAll(structCode, `;type:"`, `"`)
		structCode = strings.ReplaceAll(structCode, `;type:;primary_key`, ``)
		structCode = strings.ReplaceAll(structCode, `bson:"id" json:"id"`, `bson:"_id" json:"id"`)
	}

	return structCode, newImportPaths, nil
}

func getModelCode(data modelCodes) (string, error) {
	builder := strings.Builder{}
	err := modelTmpl.Execute(&builder, data)
	if err != nil {
		return "", err
	}

	code, err := format.Source([]byte(builder.String()))
	if err != nil {
		return "", fmt.Errorf("format.Source error: %v", err)
	}

	return string(code), nil
}

func getUpdateFieldsCode(data tmplData, isEmbed bool) (string, error) {
	_ = isEmbed

	// filter fields
	var newFields = []tmplField{}
	for _, field := range data.Fields {
		falseColumns := []string{}
		if isIgnoreFields(field.ColName, falseColumns...) || field.ColName == columnID || field.ColName == _columnID {
			continue
		}
		newFields = append(newFields, field)
	}
	data.Fields = newFields

	buf := new(bytes.Buffer)
	err := updateFieldTmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func getHandlerStructCodes(data tmplData) (string, error) {
	newFields := []tmplField{}
	for _, field := range data.Fields {
		if field.DBDriver == DBDriverMongodb { // mongodb
			if field.Name == "ID" {
				field.GoType = "string"
			}
			if "*"+field.Name == field.GoType {
				field.GoType = "*model." + field.Name
			}
			if strings.Contains(field.GoType, "[]*") {
				field.GoType = "[]*model." + strings.ReplaceAll(field.GoType, "[]*", "")
			}
		}
		newFields = append(newFields, field)
	}
	data.Fields = newFields

	postStructCode, err := tmplExecuteWithFilter(data, handlerCreateStructTmpl)
	if err != nil {
		return "", fmt.Errorf("handlerCreateStructTmpl error: %v", err)
	}

	putStructCode, err := tmplExecuteWithFilter(data, handlerUpdateStructTmpl, columnID)
	if err != nil {
		return "", fmt.Errorf("handlerUpdateStructTmpl error: %v", err)
	}

	getStructCode, err := tmplExecuteWithFilter(data, handlerDetailStructTmpl, columnID, columnCreatedAt, columnUpdatedAt)
	if err != nil {
		return "", fmt.Errorf("handlerDetailStructTmpl error: %v", err)
	}

	return postStructCode + putStructCode + getStructCode, nil
}

// customized filter fields
func tmplExecuteWithFilter(data tmplData, tmpl *template.Template, reservedColumns ...string) (string, error) {
	var newFields = []tmplField{}
	for _, field := range data.Fields {
		if isIgnoreFields(field.ColName, reservedColumns...) {
			continue
		}
		if field.DBDriver == DBDriverMongodb { // mongodb
			if strings.ToLower(field.Name) == "id" {
				field.GoType = "string"
			}
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

func getModelJSONCode(data tmplData) (string, error) {
	builder := strings.Builder{}
	err := modelJSONTmpl.Execute(&builder, data)
	if err != nil {
		return "", err
	}

	code, err := format.Source([]byte(builder.String()))
	if err != nil {
		return "", fmt.Errorf("format.Source error: %v", err)
	}

	modelJSONCode := strings.ReplaceAll(string(code), " =", ":")
	modelJSONCode = addCommaToJSON(modelJSONCode)

	return modelJSONCode, nil
}

func getProtoFileCode(data tmplData, isWebProto bool) (string, error) {
	data.Fields = goTypeToProto(data.Fields)

	builder := strings.Builder{}
	if isWebProto {
		err := protoFileForWebTmpl.Execute(&builder, data)
		if err != nil {
			return "", err
		}
	} else {
		err := protoFileTmpl.Execute(&builder, data)
		if err != nil {
			return "", err
		}
	}
	code := builder.String()

	protoMessageCreateCode, err := tmplExecuteWithFilter(data, protoMessageCreateTmpl)
	if err != nil {
		return "", fmt.Errorf("handlerCreateStructTmpl error: %v", err)
	}

	protoMessageUpdateCode, err := tmplExecuteWithFilter(data, protoMessageUpdateTmpl, columnID)
	if err != nil {
		return "", fmt.Errorf("handlerCreateStructTmpl error: %v", err)
	}
	if !isWebProto {
		protoMessageUpdateCode = strings.ReplaceAll(protoMessageUpdateCode, ` [(tagger.tags) = "uri:\"id\"" ]`, "")
	}

	protoMessageDetailCode, err := tmplExecuteWithFilter(data, protoMessageDetailTmpl, columnID, columnCreatedAt, columnUpdatedAt)
	if err != nil {
		return "", fmt.Errorf("handlerCreateStructTmpl error: %v", err)
	}

	code = strings.ReplaceAll(code, "// protoMessageCreateCode", protoMessageCreateCode)
	code = strings.ReplaceAll(code, "// protoMessageUpdateCode", protoMessageUpdateCode)
	code = strings.ReplaceAll(code, "// protoMessageDetailCode", protoMessageDetailCode)
	code = strings.ReplaceAll(code, "*time.Time", "int64")
	code = strings.ReplaceAll(code, "time.Time", "int64")
	code = adaptedDbType(data, isWebProto, code)

	return code, nil
}

const (
	createTableReplyFieldCodeMark         = "// createTableReplyFieldCode"
	deleteTableByIDRequestFieldCodeMark   = "// deleteTableByIDRequestFieldCode"
	deleteTableByIDsRequestFieldCodeMark  = "// deleteTableByIDsRequestFieldCode"
	getTableByIDRequestFieldCodeMark      = "// getTableByIDRequestFieldCode"
	getTableByIDsRequestFieldCodeMark     = "// getTableByIDsRequestFieldCode"
	listTableByLastIDRequestFieldCodeMark = "// listTableByLastIDRequestFieldCode"
)

var grpcDefaultProtoMessageFieldCodes = map[string]string{
	createTableReplyFieldCodeMark:         "uint64 id = 1;",
	deleteTableByIDRequestFieldCodeMark:   "uint64 id = 1 [(validate.rules).uint64.gt = 0];",
	deleteTableByIDsRequestFieldCodeMark:  "repeated uint64 ids = 1 [(validate.rules).repeated.min_items = 1];",
	getTableByIDRequestFieldCodeMark:      "uint64 id = 1 [(validate.rules).uint64.gt = 0];",
	getTableByIDsRequestFieldCodeMark:     "repeated uint64 ids = 1 [(validate.rules).repeated.min_items = 1];",
	listTableByLastIDRequestFieldCodeMark: "uint64 lastID = 1; // last id",
}

var webDefaultProtoMessageFieldCodes = map[string]string{
	createTableReplyFieldCodeMark:         "uint64 id = 1;",
	deleteTableByIDRequestFieldCodeMark:   `uint64 id =1 [(validate.rules).uint64.gt = 0, (tagger.tags) = "uri:\"id\""];`,
	deleteTableByIDsRequestFieldCodeMark:  "repeated uint64 ids = 1 [(validate.rules).repeated.min_items = 1];",
	getTableByIDRequestFieldCodeMark:      `uint64 id =1 [(validate.rules).uint64.gt = 0, (tagger.tags) = "uri:\"id\"" ];`,
	getTableByIDsRequestFieldCodeMark:     "repeated uint64 ids = 1 [(validate.rules).repeated.min_items = 1];",
	listTableByLastIDRequestFieldCodeMark: `uint64 lastID = 1 [(tagger.tags) = "form:\"lastID\""]; // last id`,
}

var grpcProtoMessageFieldCodes = map[string]string{
	createTableReplyFieldCodeMark:         "string id = 1;",
	deleteTableByIDRequestFieldCodeMark:   "string id = 1 [(validate.rules).string.min_len = 6];",
	deleteTableByIDsRequestFieldCodeMark:  "repeated string ids = 1 [(validate.rules).repeated.min_items = 1];",
	getTableByIDRequestFieldCodeMark:      "string id = 1 [(validate.rules).string.min_len = 6];",
	getTableByIDsRequestFieldCodeMark:     "repeated string ids = 1 [(validate.rules).repeated.min_items = 1];",
	listTableByLastIDRequestFieldCodeMark: "string lastID = 1; // last id",
}

var webProtoMessageFieldCodes = map[string]string{
	createTableReplyFieldCodeMark:         "string id = 1;",
	deleteTableByIDRequestFieldCodeMark:   `string id =1 [(validate.rules).string.min_len = 6, (tagger.tags) = "uri:\"id\""];`,
	deleteTableByIDsRequestFieldCodeMark:  "repeated string ids = 1 [(validate.rules).repeated.min_items = 1];",
	getTableByIDRequestFieldCodeMark:      `string id =1 [(validate.rules).string.min_len = 6, (tagger.tags) = "uri:\"id\"" ];`,
	getTableByIDsRequestFieldCodeMark:     "repeated string ids = 1 [(validate.rules).repeated.min_items = 1];",
	listTableByLastIDRequestFieldCodeMark: `string lastID = 1 [(tagger.tags) = "form:\"lastID\""]; // last id`,
}

func adaptedDbType(data tmplData, isWebProto bool, code string) string {
	switch data.DBDriver {
	case DBDriverMongodb: // mongodb
		if isWebProto {
			code = replaceProtoMessageFieldCode(code, webProtoMessageFieldCodes)
		} else {
			code = replaceProtoMessageFieldCode(code, grpcProtoMessageFieldCodes)
		}
	default:
		if isWebProto {
			code = replaceProtoMessageFieldCode(code, webDefaultProtoMessageFieldCodes)
		} else {
			code = replaceProtoMessageFieldCode(code, grpcDefaultProtoMessageFieldCodes)
		}
	}

	if data.ProtoSubStructs != "" {
		code += "\n" + data.ProtoSubStructs
	}

	return code
}

func replaceProtoMessageFieldCode(code string, messageFields map[string]string) string {
	for k, v := range messageFields {
		code = strings.ReplaceAll(code, k, v)
	}
	return code
}

func getServiceStructCode(data tmplData) (string, error) {
	builder := strings.Builder{}
	err := serviceStructTmpl.Execute(&builder, data)
	if err != nil {
		return "", err
	}
	code := builder.String()

	serviceCreateStructCode, err := tmplExecuteWithFilter(data, serviceCreateStructTmpl)
	if err != nil {
		return "", fmt.Errorf("handlerCreateStructTmpl error: %v", err)
	}
	serviceCreateStructCode = strings.ReplaceAll(serviceCreateStructCode, "ID:", "Id:")

	serviceUpdateStructCode, err := tmplExecuteWithFilter(data, serviceUpdateStructTmpl, columnID)
	if err != nil {
		return "", fmt.Errorf("handlerCreateStructTmpl error: %v", err)
	}
	serviceUpdateStructCode = strings.ReplaceAll(serviceUpdateStructCode, "ID:", "Id:")

	code = strings.ReplaceAll(code, "// serviceCreateStructCode", serviceCreateStructCode)
	code = strings.ReplaceAll(code, "// serviceUpdateStructCode", serviceUpdateStructCode)

	return code, nil
}

func addCommaToJSON(modelJSONCode string) string {
	r := strings.NewReader(modelJSONCode)
	buf := bufio.NewReader(r)

	lines := []string{}
	count := 0
	for {
		line, err := buf.ReadString(byte('\n'))
		if err != nil {
			break
		}
		lines = append(lines, line)
		if len(line) > 5 {
			count++
		}
	}

	out := ""
	for _, line := range lines {
		if len(line) < 5 && (strings.Contains(line, "{") || strings.Contains(line, "}")) {
			out += line
			continue
		}
		count--
		if count == 0 {
			out += line
			continue
		}
		index := bytes.IndexByte([]byte(line), '\n')
		out += line[:index] + "," + line[index:]
	}
	return out
}

// nolint
func mysqlToGoType(colTp *types.FieldType, style NullStyle) (name string, path string) {
	if style == NullInSql {
		path = "database/sql"
		switch colTp.Tp {
		case mysql.TypeTiny, mysql.TypeShort, mysql.TypeInt24, mysql.TypeLong:
			name = "sql.NullInt32"
		case mysql.TypeLonglong:
			name = "sql.NullInt64"
		case mysql.TypeFloat, mysql.TypeDouble:
			name = "sql.NullFloat64"
		case mysql.TypeString, mysql.TypeVarchar, mysql.TypeVarString,
			mysql.TypeBlob, mysql.TypeTinyBlob, mysql.TypeMediumBlob, mysql.TypeLongBlob:
			name = "sql.NullString"
		case mysql.TypeTimestamp, mysql.TypeDatetime, mysql.TypeDate:
			name = "sql.NullTime"
		case mysql.TypeDecimal, mysql.TypeNewDecimal:
			name = "sql.NullString"
		case mysql.TypeJSON:
			name = "sql.NullString"
		default:
			return "UnSupport", ""
		}
	} else {
		switch colTp.Tp {
		case mysql.TypeTiny, mysql.TypeShort, mysql.TypeInt24, mysql.TypeLong:
			if mysql.HasUnsignedFlag(colTp.Flag) {
				name = "uint"
			} else {
				name = "int"
			}
		case mysql.TypeLonglong:
			if mysql.HasUnsignedFlag(colTp.Flag) {
				name = "uint64"
			} else {
				name = "int64"
			}
		case mysql.TypeFloat, mysql.TypeDouble:
			name = "float64"
		case mysql.TypeString, mysql.TypeVarchar, mysql.TypeVarString,
			mysql.TypeBlob, mysql.TypeTinyBlob, mysql.TypeMediumBlob, mysql.TypeLongBlob:
			name = "string"
		case mysql.TypeTimestamp, mysql.TypeDatetime, mysql.TypeDate:
			path = "time" //nolint
			name = "time.Time"
		case mysql.TypeDecimal, mysql.TypeNewDecimal:
			name = "string"
		case mysql.TypeJSON:
			name = "string"
		default:
			return "UnSupport", ""
		}
		if style == NullInPointer {
			name = "*" + name
		}
	}
	return name, path
}

// nolint
func goTypeToProto(fields []tmplField) []tmplField {
	var newFields []tmplField
	for _, field := range fields {
		switch field.GoType {
		case "int":
			field.GoType = "int32"
		case "uint":
			field.GoType = "uint32"
		case "time.Time", "*time.Time":
			field.GoType = "string"
		case "float32":
			field.GoType = "float"
		case "float64":
			field.GoType = "double"
		case goTypeInts, "[]int64":
			field.GoType = "repeated int64"
		case "[]int32":
			field.GoType = "repeated int32"
		case "[]byte":
			field.GoType = "string"
		case goTypeStrings:
			field.GoType = "repeated string"
		}

		if field.DBDriver == DBDriverMongodb {
			if field.GoType[0] == '*' {
				field.GoType = field.GoType[1:]
			} else if strings.Contains(field.GoType, "[]*") {
				field.GoType = "repeated " + strings.ReplaceAll(field.GoType, "[]*", "")
			}
			if field.GoType == "[]time.Time" {
				field.GoType = "repeated string"
			}
		} else {
			if strings.ToLower(field.Name) == "id" {
				field.GoType = "uint64"
			}
		}

		newFields = append(newFields, field)
	}
	return newFields
}

func makeTagStr(tags []string) string {
	builder := strings.Builder{}
	for i := 0; i < len(tags)/2; i++ {
		builder.WriteString(tags[i*2])
		builder.WriteString(`:"`)
		builder.WriteString(tags[i*2+1])
		builder.WriteString(`" `)
	}
	if builder.Len() > 0 {
		return builder.String()[:builder.Len()-1]
	}
	return builder.String()
}

func getDefaultValue(expr ast.ExprNode) (value string) {
	if expr.GetDatum().Kind() != types.KindNull {
		value = fmt.Sprintf("%v", expr.GetDatum().GetValue())
	} else if expr.GetFlag() != ast.FlagConstant {
		if expr.GetFlag() == ast.FlagHasFunc {
			if funcExpr, ok := expr.(*ast.FuncCallExpr); ok {
				value = funcExpr.FnName.O
			}
		}
	}
	return value
}

var acronym = map[string]struct{}{
	"ID":  {},
	"IP":  {},
	"RPC": {},
}

func toCamel(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	s += "."

	n := strings.Builder{}
	n.Grow(len(s))
	temp := strings.Builder{}
	temp.Grow(len(s))
	wordFirst := true
	for _, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'
		if wordFirst && vIsLow {
			v -= 'a' - 'A'
		}

		if vIsCap || vIsLow {
			temp.WriteByte(v)
			wordFirst = false
		} else {
			isNum := v >= '0' && v <= '9'
			wordFirst = isNum || v == '_' || v == ' ' || v == '-' || v == '.'
			if temp.Len() > 0 && wordFirst {
				word := temp.String()
				upper := strings.ToUpper(word)
				if _, ok := acronym[upper]; ok {
					n.WriteString(upper)
				} else {
					n.WriteString(word)
				}
				temp.Reset()
			}
			if isNum {
				n.WriteByte(v)
			}
		}
	}
	return n.String()
}

func firstLetterToLow(str string) string {
	if len(str) == 0 {
		return str
	}

	if (str[0] >= 'A' && str[0] <= 'Z') || (str[0] >= 'a' && str[0] <= 'z') {
		return strings.ToLower(str[:1]) + str[1:]
	}

	return str
}
