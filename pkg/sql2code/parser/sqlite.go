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

	//sql = handleID(sql)
	for k, v := range sqliteToMysqlTypeMap {
		sql = strings.ReplaceAll(sql, k, v)
	}

	return sql, nil
}

//func handleID(sql string) string {
//	re := regexp.MustCompile(`id\s+INTEGER`)
//	matches := re.FindAllStringSubmatch(sql, -1)
//
//	for _, match := range matches {
//		sql = strings.ReplaceAll(sql, match[0], " id bigint unsigned")
//	}
//
//	return sql
//}
