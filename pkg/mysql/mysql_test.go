package mysql

import (
	"fmt"
	"testing"
	"time"
)

var dsn = "root:123456@(192.168.3.37:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

func TestInit(t *testing.T) {
	db, err := Init(dsn)
	if err != nil {
		t.Error(fmt.Sprintf("connect to mysql failed, err=%v, dsn=%s", err, dsn))
		return
	}

	t.Logf("%+v", db.Name())
}

func TestInitNoTLS(t *testing.T) {
	db, err := Init(
		dsn,
		//WithLog(), // 打印所有日志
		WithSlowThreshold(time.Millisecond*100), // 只打印执行时间超过100毫秒的日志
		WithEnableTrace(),                       // 开启链路跟踪
		WithMaxIdleConns(5),
		WithMaxOpenConns(50),
		WithConnMaxLifetime(time.Minute*3),
	)
	if err != nil {
		t.Error(fmt.Sprintf("connect to mysql failed, err=%v, dsn=%s", err, dsn))
		return
	}

	t.Logf("%+v", db.Name())
}
