package ggorm

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var dsn = "root:123456@(192.168.3.37:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

func TestInitMysql(t *testing.T) {
	db, err := InitMysql(dsn, WithEnableTrace())
	if err != nil {
		// ignore test error about not being able to connect to real mysql
		t.Logf(fmt.Sprintf("connect to mysql failed, err=%v, dsn=%s", err, dsn))
		return
	}
	defer CloseDB(db)

	t.Logf("%+v", db.Name())
}

func TestInitTidb(t *testing.T) {
	db, err := InitTidb(dsn)
	if err != nil {
		// ignore test error about not being able to connect to real tidb
		t.Logf(fmt.Sprintf("connect to mysql failed, err=%v, dsn=%s", err, dsn))
		return
	}
	defer CloseDB(db)

	t.Logf("%+v", db.Name())
}

func TestInitSqlite(t *testing.T) {
	dbFile := "test_sqlite.db"
	db, err := InitSqlite(dbFile)
	if err != nil {
		// ignore test error about not being able to connect to real sqlite
		t.Logf(fmt.Sprintf("connect to sqlite failed, err=%v, dbFile=%s", err, dbFile))
		return
	}
	defer CloseDB(db)

	t.Logf("%+v", db.Name())
}

func TestInitPostgresql(t *testing.T) {
	dsn = "host=192.168.3.37 user=root password=123456 dbname=account port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := InitPostgresql(dsn, WithEnableTrace())
	if err != nil {
		// ignore test error about not being able to connect to real postgresql
		t.Logf(fmt.Sprintf("connect to postgresql failed, err=%v, dsn=%s", err, dsn))
		return
	}
	defer CloseDB(db)

	t.Logf("%+v", db.Name())
}

func Test_gormConfig(t *testing.T) {
	o := defaultOptions()
	o.apply(
		WithLogging(nil),
		WithLogging(nil, 4),
		WithSlowThreshold(time.Millisecond*100),
		WithEnableTrace(),
		WithMaxIdleConns(5),
		WithMaxOpenConns(50),
		WithConnMaxLifetime(time.Minute*3),
		WithEnableForeignKey(),
		WithLogRequestIDKey("request_id"),
		WithRWSeparation([]string{
			"root:123456@(192.168.3.37:3306)/slave1",
			"root:123456@(192.168.3.37:3306)/slave2"},
			"root:123456@(192.168.3.37:3306)/master"),
		WithGormPlugin(nil),
	)

	c := gormConfig(o)
	assert.NotNil(t, c)

	err := rwSeparationPlugin(o)
	assert.NotNil(t, err)
}

type userExample struct {
	Model `gorm:"embedded"`

	Name   string `gorm:"type:varchar(40);unique_index;not null" json:"name"`
	Age    int    `gorm:"not null" json:"age"`
	Gender string `gorm:"type:varchar(10);not null" json:"gender"`
}

func TestGetTableName(t *testing.T) {
	name := GetTableName(&userExample{})
	assert.NotEmpty(t, name)

	name = GetTableName(userExample{})
	assert.NotEmpty(t, name)

	name = GetTableName("table")
	assert.Empty(t, name)
}

func TestCloseDB(t *testing.T) {
	sqlDB := new(sql.DB)
	checkInUse(sqlDB, time.Millisecond*100)
	checkInUse(sqlDB, time.Millisecond*600)
	db := new(gorm.DB)
	defer func() { recover() }()
	_ = CloseDB(db)
}

func TestCloseSqlDB(t *testing.T) {
	db := new(gorm.DB)
	defer func() { recover() }()
	CloseSQLDB(db)
}
