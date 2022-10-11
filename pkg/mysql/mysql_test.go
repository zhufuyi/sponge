package mysql

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var dsn = "root:123456@(192.168.3.37:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

func TestInit(t *testing.T) {
	db, err := Init(dsn, WithEnableTrace())
	if err != nil {
		// 忽略无法连接真实mysql的测试错误
		t.Logf(fmt.Sprintf("connect to mysql failed, err=%v, dsn=%s", err, dsn))
		return
	}

	t.Logf("%+v", db.Name())
}

func Test_gormConfig(t *testing.T) {
	o := defaultOptions()
	o.apply(
		WithLog(),
		WithSlowThreshold(time.Millisecond*100),
		WithEnableTrace(),
		WithMaxIdleConns(5),
		WithMaxOpenConns(50),
		WithConnMaxLifetime(time.Minute*3),
		WithEnableForeignKey(),
	)

	c := gormConfig(o)
	assert.NotNil(t, c)
}

func TestGetTableName(t *testing.T) {
	name := GetTableName(&userExample{})
	assert.NotEmpty(t, name)

	name = GetTableName(userExample{})
	assert.NotEmpty(t, name)

	name = GetTableName("table")
	assert.Empty(t, name)
}
