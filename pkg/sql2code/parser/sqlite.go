package parser

import (
	"strings"

	"github.com/zhufuyi/sponge/pkg/ggorm"
)

var sqliteToMysqlTypeMap = map[string]string{
	" INTEGER ": " INT ",
	" REAL ":    " FLOAT ",
	" BOOLEAN ": " TINYINT ",
	" integer ": " INT ",
	" real ":    " FLOAT ",
	" boolean ": " TINYINT ",
}

// GetSqliteTableInfo get table info from sqlite
func GetSqliteTableInfo(dbFile string, tableName string) (string, error) {
	db, err := ggorm.InitSqlite(dbFile)
	if err != nil {
		return "", err
	}
	defer closeDB(db)

	var sql string
	err = db.Raw("select sql from sqlite_master where type = ? and name = ?", "table", tableName).Scan(&sql).Error
	if err != nil {
		return "", err
	}

	for k, v := range sqliteToMysqlTypeMap {
		sql = strings.ReplaceAll(sql, k, v)
	}

	return sql, nil
}
