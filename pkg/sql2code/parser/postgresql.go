package parser

import (
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PGField postgresql field
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
	defer closeDB(db)

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

// ConvertToSQLByPgFields convert to mysql table ddl
func ConvertToSQLByPgFields(tableName string, fields []*PGField) (string, map[string]string) {
	fieldStr := ""
	pgTypeMap := make(map[string]string) // name:type
	for _, field := range fields {
		pgTypeMap[field.Name] = getType(field)
		sqlType := toMysqlType(field)
		if field.Name == "id" {
			fieldStr += fmt.Sprintf("    %s bigint unsigned primary key,\n", field.Name)
			continue
		}
		notnullStr := "not null"
		if !field.Notnull {
			notnullStr = "null"
		}
		fieldStr += fmt.Sprintf("    `%s` %s %s comment '%s',\n", field.Name, sqlType, notnullStr, field.Comment)
	}
	fieldStr = strings.TrimSuffix(fieldStr, ",\n")
	sqlStr := fmt.Sprintf("CREATE TABLE `%s` (\n%s\n);", tableName, fieldStr)
	return sqlStr, pgTypeMap
}

// nolint
func toMysqlType(field *PGField) string {
	switch field.Type {
	case "smallint", "integer", "smallserial", "serial", "int2", "int4":
		return "int"
	case "bigint", "bigserial", "int8":
		return "bigint"
	case "real":
		return "float"
	case "decimal", "numeric", "float4", "float8":
		return "decimal"
	case "double precision":
		return "double"
	case "money":
		return "varchar(30)"
	case "character", "character varying", "varchar", "char", "bpchar":
		if field.Lengthvar > 4 {
			return fmt.Sprintf("varchar(%d)", field.Lengthvar-4)
		} else {
			return "varchar(100)"
		}
	case "text":
		return "text"
	case "timestamp":
		return "timestamp"
	case "date":
		return "date"
	case "time": //nolint
		return "time" //nolint
	case "interval":
		return "year"
	case "json", "jsonb":
		return "json"
	case "boolean", "bool":
		return "bool"
	case "bit":
		return "tinyint(1)"
	}

	// unknown type convert to varchar
	field.Type = "varchar(100)"

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

func closeDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	_ = sqlDB.Close()
}
