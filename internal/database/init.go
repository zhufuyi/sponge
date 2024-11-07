// Package database provides database client initialization.
package database

import (
	"strings"
	"sync"

	"github.com/zhufuyi/sponge/pkg/sgorm"

	"github.com/zhufuyi/sponge/internal/config"
)

var (
	gdb     *sgorm.DB
	gdbOnce sync.Once

	ErrRecordNotFound = sgorm.ErrRecordNotFound
)

// todo generate initialisation database code here
// delete the templates code start

// InitDB connect database
func InitDB() {
	dbDriver := config.Get().Database.Driver
	switch strings.ToLower(dbDriver) {
	case sgorm.DBDriverMysql, sgorm.DBDriverTidb:
		gdb = InitMysql()
	case sgorm.DBDriverPostgresql:
		gdb = InitPostgresql()
	case sgorm.DBDriverSqlite:
		gdb = InitSqlite()
	default:
		panic("InitDB error, unsupported database driver: " + dbDriver)
	}
}

// delete the templates code end

// GetDB get db
func GetDB() *sgorm.DB {
	if gdb == nil {
		gdbOnce.Do(func() {
			InitDB()
		})
	}

	return gdb
}

// CloseDB close db
func CloseDB() error {
	return sgorm.CloseDB(gdb)
}
