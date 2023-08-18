package parser

import (
	"fmt"
	"github.com/blastrain/vitess-sqlparser/tidbparser/dependency/mysql"
	"github.com/blastrain/vitess-sqlparser/tidbparser/dependency/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSql(t *testing.T) {
	sql := `CREATE TABLE t_person_info (
  age INT(11) unsigned NULL,
  id BIGINT(11) PRIMARY KEY AUTO_INCREMENT NOT NULL COMMENT 'id',
  name VARCHAR(30) NOT NULL DEFAULT 'default_name' COMMENT 'name',
  created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  login_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  sex VARCHAR(2) NULL,
  num INT(11) DEFAULT 3 NULL,
  comment TEXT
  ) COMMENT="person info";`

	codes, err := ParseSQL(sql, WithTablePrefix("t_"), WithJSONTag(0))
	assert.Nil(t, err)
	for k, v := range codes {
		assert.NotEmpty(t, k)
		assert.NotEmpty(t, v)
	}

	codes, err = ParseSQL(sql, WithTablePrefix("t_"), WithJSONTag(0), WithEmbed())
	assert.Nil(t, err)
	for k, v := range codes {
		assert.NotEmpty(t, k)
		assert.NotEmpty(t, v)
	}
}

var testData = [][]string{
	{
		"CREATE TABLE information (age INT(11) NULL);",
		"Age int `gorm:\"column:age\"`", "",
	},
	{
		"CREATE TABLE information (age BIGINT(11) NULL COMMENT 'is age');",
		"Age int64 `gorm:\"column:age\"` // is age", "",
	},
	{
		"CREATE TABLE information (id BIGINT(11) PRIMARY KEY AUTO_INCREMENT);",
		"ID int64 `gorm:\"column:id;primary_key;AUTO_INCREMENT\"`", "",
	},
	{
		"CREATE TABLE information (user_ip varchar(20));",
		"UserIP string `gorm:\"column:user_ip\"`", "",
	},
	{
		"CREATE TABLE information (created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP);",
		"CreatedAt time.Time `gorm:\"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL\"`", "time",
	},
	{
		"CREATE TABLE information (num INT(11) DEFAULT 3 NULL);",
		"Num int `gorm:\"column:num;default:3\"`", "",
	},
	{
		"CREATE TABLE information (num double(5,6) DEFAULT 31.50 NULL);",
		"Num float64 `gorm:\"column:num;default:31.50\"`", "",
	},
	{
		"CREATE TABLE information (comment TEXT);",
		"Comment string `gorm:\"column:comment\"`", "",
	},
	{
		"CREATE TABLE information (comment TINYTEXT);",
		"Comment string `gorm:\"column:comment\"`", "",
	},
	{
		"CREATE TABLE information (comment LONGTEXT);",
		"Comment string `gorm:\"column:comment\"`", "",
	},
}

func TestParseSQLs(t *testing.T) {
	for i, test := range testData {
		msg := fmt.Sprintf("sql-%d", i)
		codes, err := ParseSQL(test[0], WithNoNullType())
		if !assert.NoError(t, err, msg) {
			continue
		}
		for k, v := range codes {
			if len(v) > 100 {
				v = v[:100]
			}
			t.Log(i+1, k, v)
		}
	}
}

func Test_toCamel(t *testing.T) {
	str := "user_example"
	t.Log(toCamel(str))
}

func Test_parseOption(t *testing.T) {
	opts := []Option{
		WithCharset("foo"),
		WithCollation("foo"),
		WithTablePrefix("foo"),
		WithColumnPrefix("foo"),
		WithJSONTag(1),
		WithNoNullType(),
		WithNullStyle(1),
		WithPackage("model"),
		WithGormType(),
		WithForceTableName(),
		WithEmbed(),
	}
	o := parseOption(opts)
	assert.NotNil(t, o)
}

func Test_mysqlToGoType(t *testing.T) {
	testData := []*types.FieldType{
		{Tp: uint8('n')},
		{Tp: mysql.TypeTiny},
		{Tp: mysql.TypeLonglong},
		{Tp: mysql.TypeFloat},
		{Tp: mysql.TypeString},
		{Tp: mysql.TypeTimestamp},
		{Tp: mysql.TypeDecimal},
		{Tp: mysql.TypeJSON},
	}
	var names []string
	for _, d := range testData {
		name1, _ := mysqlToGoType(d, NullInSql)
		name2, _ := mysqlToGoType(d, NullInPointer)
		names = append(names, name1, name2)
	}
	t.Log(names)
}

func Test_goTypeToProto(t *testing.T) {
	testData := []tmplField{
		{GoType: "int"},
		{GoType: "uint"},
		{GoType: "time.Time"},
	}
	v := goTypeToProto(testData)
	assert.NotNil(t, v)
}

func TestGetTableInfo(t *testing.T) {
	info, err := GetTableInfo("root:123456@(192.168.3.37:3306)/test", "user")
	t.Log(err, info)
}

func Test_initTemplate(t *testing.T) {
	defer func() { recover() }()
	modelStructTmplRaw = "{{if .foo}}"
	modelTmplRaw = "{{if .foo}}"
	updateFieldTmplRaw = "{{if .foo}}"
	handlerCreateStructTmplRaw = "{{if .foo}}"
	handlerUpdateStructTmplRaw = "{{if .foo}}"
	handlerDetailStructTmplRaw = "{{if .foo}}"
	modelJSONTmplRaw = "{{if .foo}}"
	protoFileTmplRaw = "{{if .foo}}"
	protoFileForWebTmplRaw = "{{if .foo}}"
	protoMessageCreateTmplRaw = "{{if .foo}}"
	protoMessageUpdateTmplRaw = "{{if .foo}}"
	protoMessageDetailTmplRaw = "{{if .foo}}"
	serviceCreateStructTmplRaw = "{{if .foo}}"
	serviceUpdateStructTmplRaw = "{{if .foo}}"
	serviceStructTmplRaw = "{{if .foo}}"
	initTemplate()
}
