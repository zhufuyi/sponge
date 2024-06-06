package utils

import (
	"fmt"
	"net/url"
	"strings"
)

// AdaptiveMysqlDsn adaptation of various mysql format dsn address
func AdaptiveMysqlDsn(dsn string) string {
	return strings.ReplaceAll(dsn, "mysql://", "")
}

// AdaptivePostgresqlDsn convert postgres dsn to kv string
func AdaptivePostgresqlDsn(dsn string) string {
	if strings.Count(dsn, " ") > 3 {
		return dsn
	}

	if !strings.Contains(dsn, "postgres://") {
		dsn = "postgres://" + dsn
	}

	dsn = deleteBrackets(dsn)

	u, err := url.Parse(dsn)
	if err != nil {
		panic(err)
	}

	password, _ := u.User.Password()

	if u.RawQuery == "" {
		u.RawQuery = "sslmode=disable"
	} else if u.Query().Get("sslmode") == "" {
		u.RawQuery = "sslmode=disable&" + u.RawQuery
	}
	ss := strings.Split(u.RawQuery, "&")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s",
		u.Hostname(), u.Port(), u.User.Username(), password, u.Path[1:], strings.Join(ss, " "))
}

// AdaptiveSqlite adaptive sqlite
func AdaptiveSqlite(dbFile string) string {
	// todo convert to absolute path
	return dbFile
}

// AdaptiveMongodbDsn adaptive mongodb dsn
func AdaptiveMongodbDsn(dsn string) string {
	if !strings.Contains(dsn, "mongodb://") &&
		!strings.Contains(dsn, "mongodb+srv://") {
		dsn = "mongodb://" + dsn // default scheme
	}

	return deleteBrackets(dsn)
}

func deleteBrackets(str string) string {
	start := strings.Index(str, "@(")
	end := strings.LastIndex(str, ")/")

	if start == -1 || end == -1 {
		return str
	}

	addr := str[start+2 : end]
	return strings.Replace(str, "@("+addr+")/", "@"+addr+"/", 1)
}
