package parser

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" //nolint
)

// GetMysqlTableInfo get table info from mysql
func GetMysqlTableInfo(dsn, tableName string) (string, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return "", fmt.Errorf("GetMysqlTableInfo error, %v", err)
	}
	defer db.Close() //nolint

	rows, err := db.Query("SHOW CREATE TABLE " + tableName)
	if err != nil {
		return "", fmt.Errorf("query show create table error, %v", err)
	}

	defer rows.Close() //nolint
	if !rows.Next() {
		return "", fmt.Errorf("not found found table '%s'", tableName)
	}

	var table string
	var info string
	err = rows.Scan(&table, &info)
	if err != nil {
		return "", err
	}

	return info, nil
}

// GetTableInfo get table info from mysql
// Deprecated: replaced by GetMysqlTableInfo
func GetTableInfo(dsn, tableName string) (string, error) {
	return GetMysqlTableInfo(dsn, tableName)
}
