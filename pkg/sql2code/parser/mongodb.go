package parser

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/huandu/xstrings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	mgoOptions "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/zhufuyi/sponge/pkg/mgo"
	"github.com/zhufuyi/sponge/pkg/utils"
)

const (
	goTypeOID            = "primitive.ObjectID"
	goTypeInt            = "int"
	goTypeInt64          = "int64"
	goTypeFloat64        = "float64"
	goTypeString         = "string"
	goTypeTime           = "time.Time"
	goTypeBool           = "bool"
	goTypeNil            = "nil"
	goTypeBytes          = "[]byte"
	goTypeStrings        = "[]string"
	goTypeInts           = "[]int"
	goTypeInterface      = "interface{}"
	goTypeSliceInterface = "[]interface{}"

	// SubStructKey sub struct key
	SubStructKey = "_sub_struct_"
	// ProtoSubStructKey proto sub struct key
	ProtoSubStructKey = "_proto_sub_struct_"

	oidName = "_id"
)

var mgoTypeToGo = map[bsontype.Type]string{
	bson.TypeObjectID:         goTypeOID,
	bson.TypeInt32:            goTypeInt,
	bson.TypeInt64:            goTypeInt64,
	bson.TypeDouble:           goTypeFloat64,
	bson.TypeString:           goTypeString,
	bson.TypeArray:            goTypeSliceInterface,
	bson.TypeEmbeddedDocument: goTypeInterface,
	bson.TypeTimestamp:        goTypeTime,
	bson.TypeDateTime:         goTypeTime,
	bson.TypeBoolean:          goTypeBool,
	bson.TypeNull:             goTypeNil,
	bson.TypeBinary:           goTypeBytes,
	bson.TypeUndefined:        goTypeInterface,
	bson.TypeCodeWithScope:    goTypeString,
	bson.TypeSymbol:           goTypeString,
	bson.TypeRegex:            goTypeString,
	bson.TypeDecimal128:       goTypeInterface,
	bson.TypeDBPointer:        goTypeInterface,
	bson.TypeMinKey:           goTypeInt,
	bson.TypeMaxKey:           goTypeInt,
	bson.TypeJavaScript:       goTypeString,
}

var jsonTagFormat int32 = 1 // 0: snake case, 1: camel case

// SetJSONTagSnakeCase set json tag format to snake case
func SetJSONTagSnakeCase() {
	atomic.AddInt32(&jsonTagFormat, -jsonTagFormat)
}

// SetJSONTagCamelCase set json tag format to camel case
func SetJSONTagCamelCase() {
	atomic.AddInt32(&jsonTagFormat, 1)
}

// MgoField mongo field
type MgoField struct {
	Name           string `json:"name"`
	Type           string `json:"type"`
	Comment        string `json:"comment"`
	ObjectStr      string `json:"objectStr"`
	ProtoObjectStr string `json:"protoObjectStr"`
}

// GetMongodbTableInfo get table info from mongodb
func GetMongodbTableInfo(dsn string, tableName string) ([]*MgoField, error) {
	timeout := time.Second * 5
	opts := &mgoOptions.ClientOptions{Timeout: &timeout}
	dsn = utils.AdaptiveMongodbDsn(dsn)
	db, err := mgo.Init(dsn, opts)
	if err != nil {
		return nil, err
	}

	return getMongodbTableFields(db, tableName)
}

func getMongodbTableFields(db *mongo.Database, collectionName string) ([]*MgoField, error) {
	findOpts := new(mgoOptions.FindOneOptions)
	findOpts.Sort = bson.M{oidName: -1}
	result := db.Collection(collectionName).FindOne(context.Background(), bson.M{}, findOpts)
	raw, err := result.Raw()
	if err != nil {
		return nil, err
	}

	elements, err := raw.Elements()
	if err != nil {
		return nil, err
	}

	fields := []*MgoField{}
	names := []string{}
	for _, element := range elements {
		name := element.Key()
		if name == "deleted_at" { // filter deleted_at, used for soft delete
			continue
		}
		names = append(names, name)
		t, o, p := getTypeFromMgo(name, element)
		fields = append(fields, &MgoField{
			Name:           name,
			Type:           t,
			ObjectStr:      o,
			ProtoObjectStr: p,
		})
	}

	return embedTimeField(names, fields), nil
}

func getTypeFromMgo(name string, element bson.RawElement) (goTypeStr string, goObjectStr string, protoObjectStr string) {
	v := element.Value()
	switch v.Type {
	case bson.TypeEmbeddedDocument:
		var br bson.Raw = v.Value
		es, err := br.Elements()
		if err != nil {
			return goTypeInterface, "", ""
		}
		return parseObject(name, es)

	case bson.TypeArray:
		var br bson.Raw = v.Value
		es, err := br.Elements()
		if err != nil {
			return goTypeInterface, "", ""
		}
		if len(es) == 0 {
			return goTypeInterface, "", ""
		}
		t, o, p := parseArray(name, es[0])
		return convertToSingular(t, o, p)
	}

	if goType, ok := mgoTypeToGo[v.Type]; !ok {
		return goTypeInterface, "", ""
	} else { //nolint
		return goType, "", ""
	}
}

func parseObject(name string, elements []bson.RawElement) (goTypeStr string, goObjectStr string, protoObjectStr string) {
	var goObjStr string
	var protoObjStr string
	for num, element := range elements {
		t, _, _ := getTypeFromMgo(name, element)
		k := element.Key()

		var jsonTag string
		if jsonTagFormat == 0 {
			jsonTag = xstrings.ToSnakeCase(k)
		} else {
			jsonTag = toLowerFirst(xstrings.ToCamelCase(k))
		}

		goObjStr += fmt.Sprintf("    %s %s `bson:\"%s\" json:\"%s\"`\n", xstrings.ToCamelCase(k), t, k, jsonTag)
		num++
		protoObjStr += fmt.Sprintf("  %s %s = %d;\n", convertToProtoFieldType(name, t), k, num)
	}
	return "*" + xstrings.ToCamelCase(name),
		fmt.Sprintf("type %s struct {\n%s}\n", xstrings.ToCamelCase(name), goObjStr),
		fmt.Sprintf("message %s {\n%s}\n", xstrings.ToCamelCase(name), protoObjStr)
}

func parseArray(name string, element bson.RawElement) (goTypeStr string, goObjectStr string, protoObjectStr string) {
	t, o, p := getTypeFromMgo(name, element)
	if o != "" {
		return "[]" + t, o, p
	}
	return "[]" + t, "", ""
}

func toLowerFirst(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToLower(string(str[0])) + str[1:]
}

func embedTimeField(names []string, fields []*MgoField) []*MgoField {
	isHaveCreatedAt, isHaveUpdatedAt := false, false
	for _, name := range names {
		if name == "created_at" || name == "createdAt" {
			isHaveCreatedAt = true
		}
		if name == "updated_at" || name == "updatedAt" {
			isHaveUpdatedAt = true
		}
		names = append(names, name)
	}

	var timeFields []*MgoField
	if !isHaveCreatedAt {
		timeFields = append(timeFields, &MgoField{
			Name: "created_at",
			Type: goTypeTime,
		})
	}
	if !isHaveUpdatedAt {
		timeFields = append(timeFields, &MgoField{
			Name: "updated_at",
			Type: goTypeTime,
		})
	}

	if len(timeFields) == 0 {
		return fields
	}

	return append(fields, timeFields...)
}

// ConvertToSQLByMgoFields convert to mysql table ddl
func ConvertToSQLByMgoFields(tableName string, fields []*MgoField) (string, map[string]string) {
	isHaveID := false
	fieldStr := ""
	srcMongoTypeMap := make(map[string]string) // name:type
	objectStrs := []string{}
	protoObjectStrs := []string{}

	for _, field := range fields {
		switch field.Type {
		case goTypeInterface, goTypeSliceInterface:
			srcMongoTypeMap[field.Name] = xstrings.ToCamelCase(field.Name)
		default:
			srcMongoTypeMap[field.Name] = field.Type
		}
		if field.Name == oidName {
			isHaveID = true
			srcMongoTypeMap["id"] = field.Type
			continue
		}

		fieldStr += fmt.Sprintf("    `%s` %s,\n", field.Name, convertMongoToMysqlType(field.Type))
		if field.ObjectStr != "" {
			objectStrs = append(objectStrs, field.ObjectStr)
			protoObjectStrs = append(protoObjectStrs, field.ProtoObjectStr)
		}
	}

	fieldStr = strings.TrimSuffix(fieldStr, ",\n")
	if isHaveID {
		fieldStr = "    `id` varchar(24),\n" + fieldStr + ",\n    PRIMARY KEY (id)"
	}

	if len(objectStrs) > 0 {
		srcMongoTypeMap[SubStructKey] = strings.Join(objectStrs, "\n") + "\n"
		srcMongoTypeMap[ProtoSubStructKey] = strings.Join(protoObjectStrs, "\n") + "\n"
	}

	return fmt.Sprintf("CREATE TABLE `%s` (\n%s\n);", tableName, fieldStr), srcMongoTypeMap
}

// nolint
func convertMongoToMysqlType(goType string) string {
	switch goType {
	case goTypeInt:
		return "int"
	case goTypeInt64:
		return "bigint"
	case goTypeFloat64:
		return "double"
	case goTypeString:
		return "varchar(255)"
	case goTypeTime:
		return "timestamp" //nolint
	case goTypeBool:
		return "tinyint(1)"
	case goTypeOID, goTypeNil, goTypeBytes, goTypeInterface, goTypeSliceInterface, goTypeInts, goTypeStrings:
		return "json"
	}
	return "json"
}

// nolint
func convertToProtoFieldType(name string, goType string) string {
	switch goType {
	case "int":
		return "int32"
	case "uint":
		return "uint32" //nolint
	case "time.Time":
		return "int64"
	case "float32":
		return "float"
	case "float64":
		return "double"
	case goTypeInts, "[]int64":
		return "repeated int64"
	case "[]int32":
		return "repeated int32"
	case "[]byte":
		return "string"
	case goTypeStrings:
		return "repeated string"
	}

	if strings.Contains(goType, "[]") {
		t := strings.TrimLeft(goType, "[]")
		if strings.Contains(name, t) {
			return "repeated " + t
		}
	}

	return goType
}

// MgoFieldToGoStruct convert to go struct
func MgoFieldToGoStruct(name string, fs []*MgoField) string {
	var str = ""
	var objStr string

	for _, f := range fs {
		if f.Name == oidName {
			str += "    ID primitive.ObjectID `bson:\"_id\" json:\"id\"`\n"
			continue
		}
		if f.Type == goTypeInterface || f.Type == goTypeSliceInterface {
			f.Type = xstrings.ToCamelCase(f.Name)
		}
		str += fmt.Sprintf("    %s %s `bson:\"%s\" json:\"%s\"`\n", xstrings.ToCamelCase(f.Name), f.Type, f.Name, f.Name)
		if f.ObjectStr != "" {
			objStr += f.ObjectStr + "\n"
		}
	}

	return fmt.Sprintf("type %s struct {\n%s}\n\n%s\n", xstrings.ToCamelCase(name), str, objStr)
}

func toSingular(word string) string {
	if strings.HasSuffix(word, "es") {
		if len(word) > 2 {
			return word[:len(word)-2]
		}
	} else if strings.HasSuffix(word, "s") {
		if len(word) > 1 {
			return word[:len(word)-1]
		}
	}
	return word
}

func nameToSingular(goTypeStr string, targetObjectStr string, markStr string) string {
	name := strings.ReplaceAll(goTypeStr, "[]*", "")
	prefixStr := markStr + " " + name
	l := len(prefixStr)
	if len(targetObjectStr) <= l {
		return targetObjectStr
	}

	if prefixStr == targetObjectStr[:l] {
		targetObjectStr = toSingular(prefixStr) + " " + targetObjectStr[l:]
		return targetObjectStr
	}
	return targetObjectStr
}

func convertToSingular(goTypeStr string, objectStr string, protoObjectStr string) (tStr string, oStr string, pObjStr string) {
	if !strings.Contains(goTypeStr, "[]*") || objectStr == "" {
		return goTypeStr, objectStr, protoObjectStr
	}

	objectStr = nameToSingular(goTypeStr, objectStr, "type")
	protoObjectStr = nameToSingular(goTypeStr, protoObjectStr, "message")
	goTypeStr = toSingular(goTypeStr)

	return goTypeStr, objectStr, protoObjectStr
}
