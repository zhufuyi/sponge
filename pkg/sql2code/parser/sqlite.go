package parser

import (
	"fmt"
	"strings"

	"github.com/zhufuyi/sponge/pkg/ggorm"
)

// GetSqliteTableInfo get table info from sqlite
func GetSqliteTableInfo(dbFile string, tableName string) (string, error) {
	db, err := ggorm.InitSqlite(dbFile)
	if err != nil {
		return "", err
	}
	defer closeDB(db)

	var sqliteFields SqliteFields
	sql := fmt.Sprintf("PRAGMA table_info('%s')", tableName)
	err = db.Raw(sql).Scan(&sqliteFields).Error
	if err != nil {
		return "", err
	}

	return convertToSQLBySqliteFields(tableName, sqliteFields), nil
}

// SqliteField sqlite field struct
type SqliteField struct {
	Cid          int    `gorm:"column:cid" json:"cid"`
	Name         string `gorm:"column:name" json:"name"`
	Type         string `gorm:"column:type" json:"type"`
	Notnull      int    `gorm:"column:notnull" json:"notnull"`
	DefaultValue string `gorm:"column:dflt_value" json:"dflt_value"`
	Pk           int    `gorm:"column:pk" json:"pk"`
}

var sqliteToMysqlType = map[string]string{
	"integer":       "INT",
	"text":          "TEXT",
	"real":          "FLOAT",
	"datetime":      "DATETIME",
	"blob":          "BLOB",
	"boolean":       "TINYINT",
	"numeric":       " VARCHAR(255)",
	"autoincrement": "auto_increment",
}

func (field *SqliteField) getMysqlType() string {
	sqliteType := strings.ToLower(field.Type)
	if mysqlType, ok := sqliteToMysqlType[sqliteType]; ok {
		if field.Name == "id" && sqliteType == "text" {
			return "VARCHAR(50)"
		}
		return mysqlType
	}
	return "VARCHAR(100)"
}

// SqliteFields sqlite fields
type SqliteFields []*SqliteField

func (fields SqliteFields) getPrimaryField() *SqliteField {
	var f *SqliteField
	for _, field := range fields {
		if field.Pk == 1 || field.Name == "id" {
			f = field
			return f
		}
	}
	/*
		// if no primary key, find the first xxx_id field
		if f == nil {
			for _, field := range fields {
				if strings.HasSuffix(field.Name, "_id") {
					f = field
					f.Pk = 1
					return f
				}
			}
		}

		// if no xxx_id field, find the first field
		if f == nil {
			for _, field := range fields {
				f = field
				f.Pk = 1
				return f
			}
		}
	*/
	return f
}

func convertToSQLBySqliteFields(tableName string, fields SqliteFields) string {
	if len(fields) == 0 {
		return ""
	}

	fieldStr := ""
	for _, field := range fields {
		notnullStr := "not null"
		if field.Notnull == 0 {
			notnullStr = "null"
		}
		fieldStr += fmt.Sprintf("    `%s` %s %s comment '%s',\n", field.Name, field.getMysqlType(), notnullStr, "")
	}

	primaryField := fields.getPrimaryField()
	if primaryField != nil {
		fieldStr += fmt.Sprintf("    PRIMARY KEY (`%s`)\n", primaryField.Name)
	} else {
		fieldStr = strings.TrimSuffix(fieldStr, ",\n")
	}
	return fmt.Sprintf("CREATE TABLE `%s` (\n%s\n);", tableName, fieldStr)
}
