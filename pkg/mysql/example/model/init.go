package model

import (
	"sync"

	"github.com/zhufuyi/sponge/pkg/mysql"

	"gorm.io/gorm"
)

var (
	// ErrNotFound 空记录
	ErrNotFound = gorm.ErrRecordNotFound
)

var (
	db   *gorm.DB
	once sync.Once
	dsn  string
)

// InitMysql 连接mysql
func InitMysql(addr string) {
	dsn = addr
	var err error
	db, err = mysql.Init(addr, mysql.WithLog())
	if err != nil {
		panic("config.Get() error: " + err.Error())
	}
}

// GetDB 返回db对象
func GetDB() *gorm.DB {
	if db == nil {
		once.Do(func() {
			InitMysql(dsn)
		})
	}

	return db
}

// CloseDB 关闭连接
func CloseDB() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		if sqlDB != nil {
			return sqlDB.Close()
		}
	}

	return nil
}
