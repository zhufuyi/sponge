package parser

import (
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PGField struct {
	Name      string `gorm:"column:name;" json:"name"`
	Type      string `gorm:"column:type;" json:"type"`
	Comment   string `gorm:"column:comment;" json:"comment"`
	Length    int    `gorm:"column:length;" json:"length"`
	Lengthvar int    `gorm:"column:lengthvar;" json:"lengthvar"`
	Notnull   bool   `gorm:"column:notnull;" json:"notnull"`
}

// GetPostgresqlTableInfo get table info from postgres
func GetPostgresqlTableInfo(dsn string, tableName string) ([]*PGField, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("GetPostgresqlTableInfo error: %v", err)
	}

	return getPostgresqlTableFields(db, tableName)
}

func getPostgresqlTableFields(db *gorm.DB, tableName string) ([]*PGField, error) {
	query := fmt.Sprintf(`SELECT a.attname AS name, t.typname AS type, a.attlen AS length, a.atttypmod AS lengthvar, a.attnotnull AS notnull, b.description AS comment
FROM pg_class c, pg_attribute a
    LEFT JOIN pg_description b
    ON a.attrelid = b.objoid
        AND a.attnum = b.objsubid, pg_type t
WHERE c.relname = '%s'
    AND a.attnum > 0
    AND a.attrelid = c.oid
    AND a.atttypid = t.oid
ORDER BY a.attnum;`, tableName)

	var fields []*PGField
	result := db.Raw(query).Scan(&fields)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get table fields: %v", result.Error)
	}

	return fields, nil
}

// ConvertToMysqlTable convert to mysql table ddl
func ConvertToMysqlTable(tableName string, fields []*PGField) (string, map[string]string) {
	fieldStr := ""
	pgTypeMap := make(map[string]string) // name:type
	for _, field := range fields {
		pgTypeMap[field.Name] = getType(field)
		if field.Name == "id" {
			fieldStr += fmt.Sprintf("    %s bigint unsigned primary key,\n", field.Name)
			continue
		}
		notnullStr := "not null"
		if !field.Notnull {
			notnullStr = "null"
		}
		fieldStr += fmt.Sprintf("    %s %s %s comment '%s',\n", field.Name, toMysqlType(field), notnullStr, field.Comment)
	}
	fieldStr = strings.TrimSuffix(fieldStr, ",\n")

	return fmt.Sprintf("CREATE TABLE %s (\n%s\n);", tableName, fieldStr), pgTypeMap
}

func toMysqlType(field *PGField) string {
	switch field.Type {
	// 整型
	case "smallint", "integer", "smallserial", "serial", "int2", "int4":
		return "int"
	case "bigint", "bigserial", "int8":
		return "bigint"
	// 浮点数
	case "real":
		return "float"
	case "decimal", "numeric":
		return "decimal"
	case "double precision":
		return "double"
	case "money":
		return "varchar(30)"
	// 字符串
	case "character", "character varying", "varchar", "char", "bpchar":
		if field.Lengthvar > 4 {
			return fmt.Sprintf("varchar(%d)", field.Lengthvar-4)
		}
	case "text":
		return "text"
	// 日期时间
	case "timestamp":
		return "timestamp"
	case "date":
		return "date"
	case "time": //nolint
		return "time" //nolint
	case "interval":
		return "year"
	case "boolean":
		return "tinyint(1)"
	}
	return field.Type
}

func getType(field *PGField) string {
	switch field.Type {
	case "character", "character varying", "varchar", "char", "bpchar":
		if field.Lengthvar > 4 {
			return fmt.Sprintf("varchar(%d)", field.Lengthvar-4)
		}
	}
	return field.Type
}
