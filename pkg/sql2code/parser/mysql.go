package parser

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" //nolint
)

// GetTableInfo get table info from mysql
func GetTableInfo(dsn, tableName string) (string, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return "", fmt.Errorf("connect mysql error, %v", err)
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
